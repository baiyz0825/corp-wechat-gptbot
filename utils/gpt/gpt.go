package gpt

import (
	"bufio"
	"bytes"
	"container/list"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"corp-webot/config"
	"corp-webot/xconst"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
)

var Aksk string
var httpClient *http.Client
var uuId uuid.UUID
var commonHeader http.Header

var (
	// userName = config.GetGptConf().UserName
	// password = config.GetGptConf().Passwd
	userName = ""
	password = ""
)

// 初始化gpt包
func LoadGptUtils() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		logrus.Fatal("初始化cookie客户端失败")
	}
	httpClient = &http.Client{
		Timeout: time.Second * 10,
		// 需要保证同一会话，自动管理cookie
		Jar: jar,
	}

	// 防止多次重定向
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return &xconst.GPTREDIRECTERR
	}
	httpClient.Transport = &http2.Transport{
		// 跳过证书验证
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	uuId = uuid.New()
	commonHeader = http.Header{
		"Accept-Language": []string{"zh-CN,zh;q=0.9"},
		"Accept":          []string{"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"User-Agent":      []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"},
		"Connection":      []string{"keep-alive"},
		"Content-Type":    []string{"application/json"},
	}
	// 初始化用户名密码
	userName = config.GetGptConf().UserName
	password = config.GetGptConf().Passwd
	// 初始化Aksk
	LoadAccessToken()

}

// checkV1RespCorrect 检查数据流中数据
func checkV1RespCorrect(respData *messageGPTV1Resp, traceId string) error {
	if respData == nil {
		return errors.New("检查GPT接口传递的数据响应为空，traceID：" + traceId)
	}
	if respData.Error != nil {
		return errors.New("检查GPT接口传递的数据,远程处理错误，traceID：" + traceId + "原始错误是：" + respData.Error.(string))
	}
	return nil
}

// LoadAccessToken 加载Token
func LoadAccessToken() {
	Aksk = getAccessToken()
}

// getAccessToken 获取aksk密钥
func getAccessToken() string {
	// 检查连通性 https://explorer.api.openai.com/auth/login
	link := "https://explorer.api.openai.com/"
	logrus.WithField("Link", link).Debugf("阶段1:开始获取AkskToken....")
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		logrus.WithError(err).Errorf("阶段1:创建Http请求失败")
		return ""
	}

	// Prepare Set Header
	req.Header.Set("Host", "explorer.api.openai.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")

	response, err := httpClient.Do(req)
	defer response.Body.Close()
	if err != nil {
		logrus.WithError(err).Errorf("阶段1:创建请求失败！请检查网络连接！")
		return ""
	}

	// 正确响应
	if response.StatusCode != 200 {
		logrus.WithFields(
			logrus.Fields{
				"状态码":  response.StatusCode,
				"响应数据": response,
			},
		).Errorf("阶段1:解析响应数据失败")
		return ""
	}
	return partTwo()

}

// 获取Cf盾
func partTwo() string {
	link := "https://explorer.api.openai.com/api/auth/csrf"
	logrus.WithField("Link", link).Infof("阶段2:准备获取CF盾....")
	req, err := http.NewRequest("GET", link, nil)
	// Prepare Header
	req.Header.Set("Host", "ask.openai.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://explorer.api.openai.com/auth/login")
	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.WithError(err).Errorf("阶段2:创建Http请求失败")
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && !strings.Contains(resp.Header.Get("Content-Type"), "json") {
		logrus.WithFields(logrus.Fields{
			"状态码":  resp.StatusCode,
			"响应信息": resp,
		}).Errorf("阶段2:Htpp未能正确响应")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Errorf("阶段2:解析响应体失败")
		return ""
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		logrus.WithError(err).Errorf("阶段2:发序列化body -> map失败")
		return ""
	}
	csrfToken, ok := data["csrfToken"].(string)
	if !ok {
		logrus.Errorf("cant find __Host-next-auth.csrf-token in headers")
		return ""
	}
	return partThree(csrfToken)
}

// https://explorer.api.openai.com/api/auth/signin/auth0?
func partThree(token string) string {
	link := "https://explorer.api.openai.com/api/auth/signin/auth0?"
	logrus.WithField("Link", link).Info("阶段3:开始进行模拟登陆....")
	payload := strings.NewReader(fmt.Sprintf("callbackUrl=%%2F&csrfToken=%s&json=true", token))
	req, err := http.NewRequest("POST", link, payload)
	if err != nil {
		logrus.WithError(err).Errorf("阶段1:创建Http请求失败")
		return ""
	}

	// Prepare Header
	req.Header.Add("Authority", "explorer.api.openai.com")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Add("Origin", "https://explorer.api.openai.com")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Referer", "https://explorer.api.openai.com")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	resp, err := httpClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		logrus.WithError(err).Errorf("阶段3:创建请求失败！请检查网络连接！")
		return ""
	}

	// 校验获取信息
	contentType := resp.Header.Get("Content-Type")
	if resp.StatusCode != 200 && !strings.Contains(contentType, "json") {
		logrus.WithFields(logrus.Fields{
			"状态码":  resp.StatusCode,
			"响应信息": resp,
		}).Errorf("阶段2:Htpp可能未能正确响应")
		return ""
	}
	var jsonResponse struct {
		URL string `json:"url"`
	}
	data, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(data, &jsonResponse)
	if err != nil {
		logrus.WithError(err).Errorf("阶段3:解析响应体失败")
		return ""
	}

	// 判断是否被限制登陆
	if jsonResponse.URL == "https://explorer.api.openai.com/api/auth/error?error=OAuthSignin" || strings.Contains(jsonResponse.URL, "error") {
		logrus.WithField("地址", jsonResponse.URL).Errorf("阶段3:您已经被限制登陆")
		return ""
	}
	return partFore(jsonResponse.URL)
}

// https://auth0.openai.com/authorize?client_id
func partFore(link string) string {
	logrus.WithField("Link", link).Info("阶段4:获取登陆表单....")
	req, err := http.NewRequest("GET", link, nil)
	// Prepare Header
	req.Header.Set("Authority", "auth0.openai.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil && !errors.Is(err, &xconst.GPTREDIRECTERR) {
		logrus.WithError(err).Errorf("阶段4:出现非自定义重定向一次限制错误，可能为http请求失败！")
		return ""
	}

	// 判断是否为成功重定向
	if resp.StatusCode != 302 {
		logrus.WithFields(logrus.Fields{
			"状态码":  resp.StatusCode,
			"响应信息": resp,
		}).Errorf("阶段4:Htpp可能不是302重定向，无法获取下一阶段登陆信息")
		return ""
	}
	// 匹配状态信息
	location := resp.Header.Get("Location")
	state := regexp.MustCompile(`state=(.*)`).FindStringSubmatch(location)[0]
	state = strings.Split(state, "=")[1]
	return partFive(state)
}

// https://auth0.openai.com/u/login/identifier?state
func partFive(state string) string {
	link := fmt.Sprintf("https://auth0.openai.com/u/login/identifier?state=%s", state)
	logrus.WithField("Link", link).Info("阶段5:登陆表单获取前回调....")

	req, err := http.NewRequest("GET", link, nil)
	req.Header.Set("Host", "auth0.openai.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://explorer.api.openai.com/")
	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		logrus.WithError(err).Errorf("阶段5:发送获取实际表单信息请求出错！")
		return ""
	}
	// 成功响应
	if resp.StatusCode != 200 {
		logrus.WithFields(logrus.Fields{
			"状态码":  resp.StatusCode,
			"响应信息": resp,
		}).Errorf("阶段5:Http响应非200，请检查网络")
		return ""
	}
	return partSix(state)
}

// https://auth0.openai.com/u/login/identifier?state
func partSix(state string) string {
	link := fmt.Sprintf("https://auth0.openai.com/u/login/identifier?state=%s", state)
	logrus.WithField("Link", link).Info("阶段6:开始进行登陆..输入用户名中....")
	// 设置登陆参数
	usernameEncode := url.QueryEscape(userName)
	payload := fmt.Sprintf("state=%s&username=%s&action=default&js-available=false&webauthn-available=true&is-brave=false&webauthn-platform-available=true&action=default", state, usernameEncode)
	// 构建表单post请求 - 请求密码框
	req, err := http.NewRequest("POST", link, bytes.NewBufferString(payload))
	req.Header.Add("Authority", "explorer.api.openai.com")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Referer", fmt.Sprintf("https://auth0.openai.com/u/login/identifier?state=%s", state))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil && !errors.Is(err, &xconst.GPTREDIRECTERR) {
		logrus.WithError(err).Errorf("阶段6:出现非自定义重定向一次限制错误，可能为http请求失败！")
		return ""
	}

	if resp.StatusCode != 302 {
		logrus.WithFields(logrus.Fields{
			"状态码":  resp.StatusCode,
			"响应信息": resp,
		}).Errorf("阶段6:Htpp可能不是302重定向，无法获取下一阶段登陆信息")
		return ""
	}
	return partSeven(state)
}

// https://auth0.openai.com/u/login/password?state
func partSeven(state string) string {
	link := fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", state)
	logrus.WithField("Link", link).Info("阶段7:开始进行登陆..输入密码中....")
	// 设置请求参数
	usernameEncode := url.QueryEscape(userName)
	passwd := url.QueryEscape(password)

	payload := fmt.Sprintf("state=%s&username=%s&password=%s&action=default", state, usernameEncode, passwd)
	// 发起登陆请求
	req, err := http.NewRequest("POST", link, bytes.NewBufferString(payload))
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", state))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	defer resp.Body.Close()

	if err != nil && !errors.Is(err, &xconst.GPTREDIRECTERR) {
		logrus.WithError(err).Errorf("阶段7:出现非自定义重定向一次限制错误，可能为http请求失败！")
		return ""
	}
	if resp.StatusCode != 302 {
		logrus.WithFields(logrus.Fields{
			"状态码":  resp.StatusCode,
			"响应信息": resp,
		}).Errorf("阶段4:Htpp可能不是302重定向，无法获取下一阶段登陆信息")
		return ""

	}
	location := resp.Header.Get("Location")
	newState := regexp.MustCompile(`state=(.*)`).FindStringSubmatch(location)[0]
	newState = strings.Split(newState, "=")[1]
	return partEight(state, newState)
}

// https://auth0.openai.com/authorize/resume?state=%s
func partEight(oldState, newState string) string {
	link := fmt.Sprintf("https://auth0.openai.com/authorize/resume?state=%s", newState)
	logrus.WithField("Link", link).Info("阶段8:获取登陆表单....")
	req, err := http.NewRequest("GET", link, nil)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", oldState))

	resp, err := httpClient.Do(req)
	defer resp.Body.Close()

	// 不出错 且不是302错误
	if err != nil && !errors.Is(err, &xconst.GPTREDIRECTERR) {
		logrus.WithError(err).Errorf("阶段8:出现非自定义重定向一次限制错误，可能为http请求失败！")
		return ""
	}
	if resp.StatusCode != 302 {
		logrus.WithFields(logrus.Fields{
			"状态码":  resp.StatusCode,
			"响应信息": resp,
		}).Errorf("阶段8:Htpp可能不是302重定向，无法获取下一阶段登陆信息")
		return ""
	}
	// 获取回调地址
	return partNight(resp.Header.Get("Location"))
}

func partNight(urlCallBack string) string {
	logrus.WithField("Link", urlCallBack).Info("阶段9:进行获取token回调....")

	req, _ := http.NewRequest("GET", urlCallBack, nil)
	// Prepare Header
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Authority", "explorer.api.openai.com")

	resp, err := httpClient.Do(req)
	defer resp.Body.Close()

	// 最后一个请求有错误但是不是期待的重定向错误
	if err != nil && !errors.Is(err, &xconst.GPTREDIRECTERR) {
		logrus.WithError(err).Errorf("阶段9:发生非重定向错误")
		return ""
	}
	cookies := resp.Cookies()
	foundSessionToken := ""
	for _, cookie := range cookies {
		if cookie.Name == "__Secure-next-auth.session-token" {
			foundSessionToken = cookie.Value
			break
		}
	}
	if foundSessionToken == "" {
		logrus.WithField("Cookie", cookies).Errorf("阶段9:无法获取对应的Token")
		return ""
	}
	return getAkskFromReq(foundSessionToken)

}

func getAkskFromReq(sessionToken string) string {
	link := "https://explorer.api.openai.com/api/auth/session"
	logrus.WithField("Link", link).Info("最后阶段10:开始使用获取到的session-Token进行检查,并换取Token....")

	req, err := http.NewRequest("GET", link, nil)
	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		logrus.WithError(err).WithField("上一步获取的Token", sessionToken).Errorf("阶段10:http请求失败！")
		return ""
	}

	// 检查是否成功
	if resp.StatusCode != http.StatusOK {
		logrus.WithField("response", resp).Errorf("阶段10:http请求响应状态码失败！")
		return ""
	}

	// 反序列化保存Token
	var data struct {
		AccessToken string `json:"accessToken"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		logrus.WithError(err).Errorf("阶段10:反序列化解析AccessToken失败")
		return ""
	}

	logrus.WithField("LoginToken", data.AccessToken).Debugf("成功获取到Token")
	return data.AccessToken
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
func NewGptHelperV1(accessToken, traceID string) V1GPTHelper {
	return V1GPTHelper{
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

// GetTraceIDFromHelper  返回context traceID 不能修改
func (h *V1GPTHelper) GetTraceIDFromHelper() string {
	traceId, _ := h.ctx.Value("traceId").(string)
	return traceId
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
func (h *V1GPTHelper) SendAndGetMessageToGPTV1(client http.Client, message string) ([]byte, error) {
	// 拼接请求参数
	dataLoadStruct := messageGPTV1Req{
		Action: V1ApiAction,
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
	request.Header.Set("Authorization", "Bearer "+h.AccessToken)
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
