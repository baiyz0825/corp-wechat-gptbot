package services

import (
	"strings"
	"time"

	xcache "github.com/baiyz0825/corp-webot/cache"
	"github.com/baiyz0825/corp-webot/config"
	"github.com/baiyz0825/corp-webot/model"
	"github.com/baiyz0825/corp-webot/to"
	"github.com/baiyz0825/corp-webot/utils/openaiutils"
	"github.com/baiyz0825/corp-webot/utils/wecom"
	"github.com/sirupsen/logrus"
)

// DoChat 进行解析发送openapi -> 返回微信
func DoChat(userData to.MsgContent) bool {
	// TODO 进行command分割
	// 检查是否请求过相同内容 不存在调用openai
	respOpenAI := CompareCacheAndGetFromApi(userData)
	// 发送到微信
	mode := config.GetSystemConf().MsgMode
	switch mode {
	case "markdown":
		return SendByMarkdown(userData, respOpenAI)
	case "text":
		return SendByText(userData, respOpenAI)
	default:
		return false
	}

}

// SendByMarkdown 使用markdown发送
func SendByMarkdown(userData to.MsgContent, respOpenAI string) bool {
	// TODO: 考虑是否分片长消息分割
	resp := wecom.SendMarkdownToUSer(userData, respOpenAI)
	if resp.ResponseWork.ErrCode != 0 {
		logrus.WithField("resp:", resp).Error("企业微信助手发送失败")
		return false
	}
	return true
}

// SendByText 使用大文本发送
func SendByText(userData to.MsgContent, respOpenAI string) bool {
	// 按行分割代码
	lines := strings.Split(respOpenAI, "\n")
	// 临时消息缓冲
	var content string
	for i := 0; i < len(lines); i++ {
		// >2000 发送前一部分，清空重来
		if len(content)+len(lines[i]) > 2000 {
			resp := wecom.SendTextToUSer(userData, content)
			if resp.ResponseWork.ErrCode != 0 {
				logrus.WithField("resp:", resp).Error("企业微信助手发送分片失败")
			}
			// TODO: 重试机制
			// 清空缓存
			content = ""
		}
		// 拼接下一行
		content += lines[i] + "\n"
	}
	// 最后一个
	if content != "" {
		resp := wecom.SendTextToUSer(userData, content)
		if resp.ResponseWork.ErrCode != 0 {
			logrus.WithField("resp:", resp).Error("企业微信助手发送分片失败")
		}
	}
	return true
}

// CheckCacheUserEchoReq 检查缓存中是否存在数据 用户多少秒允许发送一次请求
func CheckCacheUserEchoReq(data to.MsgContent) bool {
	cacheKey := data.ToUsername + ":" + data.Msgid
	_, find := xcache.GetCacheDb().Get(cacheKey)
	if find {
		return true
	}
	xcache.SetDataToCache(cacheKey, "", time.Second*4)
	return false
}

// CompareCacheAndGetFromApi  比较是否重复请求查询一个内容
func CompareCacheAndGetFromApi(data to.MsgContent) string {
	// 获取上下文缓存
	var msgContext model.MessageContext
	cache := xcache.GetDataFromCache(xcache.GetUserCacheKey(data.ToUsername))
	if cache != nil {
		context, ok := cache.(model.MessageContext)
		if !ok {
			logrus.WithField("error", "上下文断言失败").
				WithField("userID", data.ToUsername).
				Errorf("用户上下文数据断言失败")
			return "稍后再试试，小助理开小差了"
		}
		msgContext = context
	} else {
		msgContext = CreateNewContextWithSysPrompt(data.ToUsername)
	}
	// 植入新的聊天内容user
	newMsg := model.NewUserMsg(data.Content)
	msgContext.Context = append(msgContext.Context, newMsg)
	// 请求openAi
	respOpenAI := openaiutils.SendReqAndGetResp(msgContext.Context)
	// 存储新的上下文内容
	assistantMsg := model.NewAssistantMsg(respOpenAI)
	msgContext.Context = append(msgContext.Context, assistantMsg)
	// update cache cache full -> delete
	if msgContext.Full() {
		xcache.GetCacheDb().Delete(msgContext.Key)
	}
	xcache.GetCacheDb().Set(msgContext.Key, msgContext, config.GetGptConf().ContextExpireTime*time.Minute)
	return respOpenAI
}

// CreateNewContextWithSysPrompt
// @Description: 创建包含提示词的prompt
// @param ketFactor
// @return model.MessageContext
func CreateNewContextWithSysPrompt(ketFactor string) model.MessageContext {
	// 查询db获取用户sysPrompt
	msg := "请使用中文和我对话"
	// 获取context
	context := model.NewUserChatContext(ketFactor)
	sysPrompt := model.NewSystemMsg(msg)
	context.Context = append(context.Context, sysPrompt)
	return context
}
