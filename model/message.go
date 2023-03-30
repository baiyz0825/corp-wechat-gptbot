package model

import (
	xcache "github.com/baiyz0825/corp-webot/cache"
	"github.com/baiyz0825/corp-webot/config"
	"github.com/sashabaranov/go-openai"
)

// MessageContext 用户上下文对象
type MessageContext struct {
	Key     string
	Context []openai.ChatCompletionMessage
}

// NewUserMsg
// @Description: 创建用户聊天信息
// @param msg
// @return *MessageContext
func NewUserMsg(msg string) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    "user",
		Content: msg,
	}
}

// NewAssistantMsg
// @Description: 创建系统回答
// @param msg
// @return *openai.ChatCompletionMessage
func NewAssistantMsg(msg string) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    "assistant",
		Content: msg,
	}
}

// NewSystemMsg
// @Description: 创建prompt提示词
// @param msg
// @return *openai.ChatCompletionMessage
func NewSystemMsg(msg string) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    "assistant",
		Content: msg,
	}
}

// NewUserChatContext
// @Description: 创建新的上下文对象
// @param keyFactor
// @return MessageContext
func NewUserChatContext(keyFactor string) MessageContext {
	return MessageContext{
		Key:     xcache.GetUserCacheKey(keyFactor),
		Context: make([]openai.ChatCompletionMessage, config.GetGptConf().ContextNumber/2, config.GetGptConf().ContextNumber),
	}
}

// Full
// @Description: context消息是否已满
// @receiver receiver
// @return bool
func (receiver MessageContext) Full() bool {
	return len(receiver.Context) == config.GetGptConf().ContextNumber*2+1
}
