package impl

import (
	"encoding/json"
	"time"

	xcache "github.com/baiyz0825/corp-webot/cache"
	"github.com/baiyz0825/corp-webot/config"
	"github.com/baiyz0825/corp-webot/dao"
	"github.com/baiyz0825/corp-webot/model"
	"github.com/baiyz0825/corp-webot/services"
	"github.com/baiyz0825/corp-webot/to"
	"github.com/baiyz0825/corp-webot/utils/openaiutils"
	"github.com/baiyz0825/corp-webot/utils/xlog"
	"github.com/baiyz0825/corp-webot/xconst"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

type GPTChatCommand struct {
	Command string
}

func NewGPTCommand() *GPTChatCommand {
	return &GPTChatCommand{}
}

// Exec
// @Description: 进行解析发送openapi -> 返回微信
// @receiver c
// @param userData
// @return bool
func (c GPTChatCommand) Exec(userData to.MsgContent) bool {
	// 检查是否请求过相同内容 不存在调用openai
	respOpenAI := CompareCacheAndGetFromApi(userData)
	// 发送到微信
	mode := config.GetSystemConf().MsgMode
	switch mode {
	case "markdown":
		return services.SendToWxByMarkdown(userData, respOpenAI)
	case "text":
		return services.SendToWxByText(userData, respOpenAI)
	default:
		return false
	}
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
			return xconst.AI_DEFAULT_MSG
		}
		msgContext = context
	} else {
		context := CreateNewContextWithSysPrompt(data.ToUsername)
		if context == nil {
			return xconst.AI_DEFAULT_MSG
		}
		msgContext = *context
	}
	// 植入新的聊天内容user
	newMsg := model.NewUserMsg(data.Content)
	msgContext.Context = append(msgContext.Context, newMsg)
	// 请求openAi
	respOpenAI := openaiutils.SendReqAndGetTextResp(msgContext.Context)
	// 存储新的上下文内容
	assistantMsg := model.NewAssistantMsg(respOpenAI)
	msgContext.Context = append(msgContext.Context, assistantMsg)
	// update cache cache full -> delete
	if msgContext.Full() {
		// 入库
		msgContextJson, err := json.Marshal(msgContext)
		if err != nil {
			xlog.Log.WithError(err).WithField("反序列化数据是", msgContextJson).
				WithField("用户是:", data.ToUsername).
				Error("系统序列化错误")
		}
		err = dao.InsertUserContext(data.ToUsername, string(msgContextJson), dao.DB)
		if err != nil {
			xlog.Log.WithError(err).WithField("插入数据是:", string(msgContextJson)).
				WithField("用户是:", data.ToUsername).
				Error("保存过期缓存中的用户上下文数据->db错误")
			return xconst.AI_DEFAULT_MSG
		}
		// 删除缓存
		xcache.GetCacheDb().Delete(msgContext.Key)
	}
	// 设置新的上下文
	xcache.GetCacheDb().Set(msgContext.Key, msgContext, config.GetGptConf().ContextExpireTime*time.Minute)
	return respOpenAI
}

// CreateNewContextWithSysPrompt
// @Description: 创建包含提示词的prompt
// @param ketFactor
// @return model.MessageContext
func CreateNewContextWithSysPrompt(userID string) *model.MessageContext {
	// 查询db获取用户sysPrompt
	user, err := dao.GetUser(userID, dao.DB)
	if err != nil {
		xlog.Log.WithError(err).Error("查询用户sysPrompt存在错误:")
		return nil
	}
	var sysPrompt openai.ChatCompletionMessage
	if len(user.SysPrompt) != 0 {
		sysPrompt = model.NewSystemMsg(user.SysPrompt)
	} else {
		sysPrompt = model.NewSystemMsg(xconst.PROMPT_DEFAULT)
	}
	// 创建context
	context := model.NewUserChatContext(xcache.GetUserCacheKey(userID))
	context.Context = append(context.Context, sysPrompt)
	return &context
}
