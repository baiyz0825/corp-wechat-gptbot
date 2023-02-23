package wxcallback

import (
	"net/http"
	"strconv"

	"corp-webot/services/command"
	"corp-webot/services/events"
	"corp-webot/utils/wecom"
	"github.com/ArtisanCloud/PowerLibs/v3/http/helper"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/contract"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type QueryParams struct {
	MsgSignature string `form:"msg_signature"`
	TimeStamp    string `form:"timestamp"`
	Nonce        string `form:"nonce"`
	EchoStr      string `form:"echostr"`
}

// VerifyCallBack 回调验证
func VerifyCallBack(c *gin.Context) {
	var q QueryParams
	if err := c.Bind(&q); err != nil {
		logrus.Errorf("绑定回调Query错误：%v", err)
	}
	echoStr, cryptErr := wecom.WxBizMsgCryptHelper.VerifyURL(q.MsgSignature, q.TimeStamp, q.Nonce, q.EchoStr)
	if cryptErr != nil {
		logrus.Errorf("验证Url出错：%v", cryptErr)
	}
	logrus.Info("解析的回调字符为：", string(echoStr))
	c.Writer.Write(echoStr)
}

// RealMsgCallBack 实际处理用户消息
func RealMsgCallBack(c *gin.Context) {
	response, err := wecom.WeComApp.Server.Notify(c.Request, func(event contract.EventInterface) interface{} {
		// 所有包含的结构体请参考： https://github.com/ArtisanCloud/PowerWeChat/tree/master/src/work/server/handlers/models
		switch event.GetMsgType() {
		case "text":
			logrus.Debug("接受到来自用户： " + event.GetFromUserName() + "的文本内容消息")
			// 开始分发处理
			go func(event contract.EventInterface) {
				// 文本消息分发函数。异步处理
				messageText := models.MessageText{}
				err := event.ReadMessage(&messageText)
				if err != nil {
					logrus.WithError(err).Errorf("处理文本消息过程中，出现读取回掉消息内容出错")
					return
				}
				// 截取命令
				commandFlag := messageText.Content[0:5]
				realMsg := messageText.Content[6:]
				// 获取响应的消息处理函数
				commandFunc, err := command.GetCommandFunc(commandFlag)
				if err != nil {
					return
				}
				agentId, err := strconv.Atoi(messageText.AgentID)
				if err != nil {
					return
				}
				// 处理消息
				commandFunc.ExecCommand(
					command.NewCommandData(event.GetToUserName(), event.GetFromUserName(), realMsg, commandFlag, agentId),
					c,
				)
			}(event)
		case "image":
			logrus.Debug("接受到来自用户： " + event.GetFromUserName() + "的图片内容消息")
		case "event":
			logrus.Debug("接受到来自用户： " + event.GetFromUserName() + "的事件消息")
			go func(event contract.EventInterface) {
				eventType := event.GetEvent()
				// 分发事件处理器
				commandFunc, err := events.GetCommandFunc(eventType)
				if err != nil {
					return
				}
				commandFunc.DealEvent(event, c)
			}(event)
		default:
			logrus.Info("接受到来自用户： " + event.GetFromUserName() + "不支持的消息类型请求（" + event.GetMsgType() + ")")
			return "我还没有学会这个功能呢/::$"
		}
		// 直接回复用户
		return "正在处理...请稍等/:,@f"
	})
	// 消息处理异常
	if err != nil {
		logrus.WithError(err).Errorf("处理用户消息错误")
		c.String(http.StatusOK, "诶呀，我这里除小差了，稍后在试试吧")
	}
	// 回送正确响应
	err = helper.HttpResponseSend(response, c.Writer)
	if err != nil {
		logrus.WithError(err).Errorf("响应用户消息错误")
	}
}
