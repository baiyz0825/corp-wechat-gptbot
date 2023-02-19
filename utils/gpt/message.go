package gpt

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var uuId uuid.UUID
var commonHeader http.Header

// init 初始化包
func init() {
	uuId = uuid.New()
	commonHeader = http.Header{
		"Accept-Language": []string{"zh-CN,zh;q=0.9"},
		"Accept":          []string{"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"User-Agent":      []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"},
		"Connection":      []string{"keep-alive"},
		"Content-Type":    []string{"application/json"},
	}
}

// messageGPTV1Req GPT APIV1模型Resp
type messageGPTV1Req struct {
	Action          string          `json:"action"`
	Messages        []messagesV1Req `json:"messages"`
	ConversationID  string          `json:"conversation_id"`
	ParentMessageID string          `json:"parent_message_id"`
	Model           string          `json:"model"`
	Paid            string          `json:"paid"`
	Stream          bool            `json:"stream"`
}
type contentV1Req struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}
type messagesV1Req struct {
	ID      string       `json:"id"`
	Role    string       `json:"role"`
	Content contentV1Req `json:"content"`
}

// messageGPTV1Resp GPT APIV1模型Resp
type messageGPTV1Resp struct {
	Message        messageV1Resp `json:"message"`
	ConversationID string        `json:"conversation_id"`
	Error          interface{}   `json:"error"`
}
type contentV1Resp struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}
type finishDetailsV1Resp struct {
	Type string `json:"type"`
	Stop string `json:"stop"`
}
type metadataResp struct {
	FinishDetails finishDetailsV1Resp `json:"finish_details"`
	MessageType   string              `json:"message_type"`
	ModelSlug     string              `json:"model_slug"`
}
type messageV1Resp struct {
	ID         string        `json:"id"`
	Role       string        `json:"role"`
	User       interface{}   `json:"user"`
	CreateTime interface{}   `json:"create_time"`
	UpdateTime interface{}   `json:"update_time"`
	Content    contentV1Resp `json:"content"`
	EndTurn    interface{}   `json:"end_turn"`
	Weight     int           `json:"weight"`
	Metadata   metadataResp  `json:"metadata"`
	Recipient  string        `json:"recipient"`
}

// V1GPTHelper GPTV!助手
type V1GPTHelper struct {
	ConVersionOneID string
	ParentMsgId     *list.List
	AccessToken     string
	ctx             context.Context
}

// NewGptHelperV1 返回一个V1 Gpt实例
func NewGptHelperV1(accessToken, traceID string) *V1GPTHelper {
	return &V1GPTHelper{
		ConVersionOneID: "",
		ParentMsgId:     list.New(),
		AccessToken:     accessToken,
		ctx:             context.WithValue(context.Background(), "traceId", traceID),
	}
}

// GetContextFromHelper 返回context 不能修改
func (h *V1GPTHelper) GetContextFromHelper() context.Context {
	return h.ctx
}

// GetMessagesFromGptV1 从数据流中解析消息
func (h *V1GPTHelper) getMessagesLineFromGptV1(resp *http.Response, traceId string, outMessageChannel chan messageGPTV1Resp, wg *sync.WaitGroup) {
	// 先关闭管道，在关闭wg
	defer close(outMessageChannel)
	defer wg.Done()
	// 初始化接受
	respData := messageGPTV1Resp{}
	// 检查http响应
	if resp == nil {
		logrus.WithField("traceId", traceId).Errorf("传递的响应为空！")
		return
	}
	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(
			logrus.Fields{
				"traceId":   traceId,
				"http-code": resp.StatusCode,
				"response":  resp,
			},
		).Errorf("响应状态码错误")
		return
	}
	// 逐行开始读取数据
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		// 截取有效数据 data: 真实数据体
		if strings.Contains(line, "data: ") {
			line = line[6:]
		}
		// 结束标记 跳出位置
		if line == "data: [DONE]" {
			break
		}
		// json unmarshal
		err := json.Unmarshal([]byte(line), &respData)
		if err != nil {
			continue
		}
		// 检查结构体完整性
		if err := checkV1RespCorrect(&respData, traceId); err != nil {
			logrus.WithError(err).Errorf("解析检查所接受的行数据结构体失败")
			continue
		}
		outMessageChannel <- respData
	}
	close(outMessageChannel)
}

// SendAndGetMessageToGPTV1 使用指定的http client发送消息并将传回的消息进行保存
func (h *V1GPTHelper) SendAndGetMessageToGPTV1(client http.Client, accessToken, message string) ([]byte, error) {
	// 拼接请求参数
	dataLoadStruct := messageGPTV1Req{
		Action: "",
		Messages: []messagesV1Req{
			{
				ID:   uuId.String(),
				Role: V1ApiRole,
				Content: contentV1Req{
					ContentType: V1ApiContentType,
					Parts:       []string{message},
				},
			},
		},
		ConversationID:  h.ConVersionOneID,
		ParentMessageID: "",
		Model:           V1ApiModel,
		Paid:            V1ApiPaid,
		Stream:          false,
	}
	// 检查是否存在上下文与存在的对话,修改请求对话上下文
	if h.ParentMsgId != nil && h.ParentMsgId.Len() > 0 {
		parentMessageID, ok := h.ParentMsgId.Back().Value.(string)
		if ok {
			dataLoadStruct.ParentMessageID = parentMessageID
		}
	}
	traceId, _ := h.ctx.Value("traceId").(string)
	jsonData, err := json.Marshal(dataLoadStruct)
	if err != nil {
		logrus.WithField("traceId", traceId).Errorf("构建Json请求体序列化错误：%v", err)
		return nil, errors.Wrap(err, "错误的traceId是："+traceId)
	}
	// 创建Http请求
	request, err := http.NewRequest("POST", V1GPTBaseUrl+"/api/conversation", bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.WithField("traceId", traceId).Errorf("构建Http请求错误：%v", err)
		return nil, errors.Wrap(err, "错误的traceId是："+traceId)
	}
	request.Header = commonHeader
	// 设置accessToken
	request.Header.Set("Authorization", "Bearer "+accessToken)
	// 发送请求
	resp, err := client.Do(request)
	if err != nil {
		logrus.WithField("traceId", traceId).Errorf("Http请求错误：%v", err)
		return nil, errors.Wrap(err, "错误的traceId是："+traceId)
	}
	// 处理响应结果
	var wg sync.WaitGroup
	// 设置读取的存储缓冲区
	messageBuffer := bytes.NewBuffer(make([]byte, 0, 4096))
	wg.Add(1)
	outMessageChannel := make(chan messageGPTV1Resp, 5)
	// 获取每一行的实际数据(预处理)
	h.getMessagesLineFromGptV1(resp, traceId, outMessageChannel, &wg)
	// 刷写存储信息
	for line := range outMessageChannel {
		h.ConVersionOneID = line.ConversationID
		// 存储当前对话响应（每次对话中，最后一次数据是最完整的，因此保存最后一次上下文,每次消息的id与对话id相等）
		h.ParentMsgId.PushBack(line.Message.ID)
		// 拼接message中的所有part
		messageBuffer.Write([]byte(strings.Join(line.Message.Content.Parts, "")))
	}
	wg.Wait()
	// 处理结束，返回缓冲区结果
	return messageBuffer.Bytes(), nil
}
