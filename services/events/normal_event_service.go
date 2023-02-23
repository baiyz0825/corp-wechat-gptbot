package events

import (
	"strconv"

	"corp-webot/utils/wecom"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/contract"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type NorMalEventService struct {
}

func (n *NorMalEventService) DealEvent(event contract.EventInterface, ctx context.Context) {
	switch event.GetEvent() {
	case models.CALLBACK_EVENT_ENTER_AGENT:
		go func() {
			data := models.EventEnterAgent{}
			err := event.ReadMessage(&data)
			if err != nil {
				logrus.WithError(err).Errorf("读取用户进入应用事件转化错误")
				return
			}
			agentId, err := strconv.Atoi(data.AgentID)
			if err != nil {
				return
			}
			messages := &request.RequestMessageSendMarkdown{
				RequestMessageSend: request.RequestMessageSend{
					ToUser:                 data.FromUserName,
					ToParty:                "",
					ToTag:                  "",
					MsgType:                "markdown",
					AgentID:                agentId,
					EnableDuplicateCheck:   0,
					DuplicateCheckInterval: 1800,
				},
				Markdown: &request.RequestMarkdown{
					Content: `
				### 欢迎进入你的智能助手
				>支持的功能为
				- @gpt：与ChatGPT聊天
				- @hp: 帮助信息
				- @暂未支持
			`,
				},
			}
			// 发送微信消息
			_, err = wecom.WeComApp.Message.SendMarkdown(ctx, messages)
			if err != nil {
				logrus.WithError(err).WithFields(
					logrus.Fields{
						"来自:": data.FromUserName,
						"应用":  agentId,
					},
				).Errorf("发送Markdown消息失败")

			}
		}()
	default:
		logrus.Debug("不支持的事件类型")
		return
	}
}
