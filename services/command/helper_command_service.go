package command

import (
	"corp-webot/utils/wecom"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type HelperCommandService struct {
}

// ExecCommand 处理帮助菜单命令
func (gpt *HelperCommandService) ExecCommand(data *CommandData, ctx context.Context) {
	// 获取用户请求中的消息
	if data != nil {
		logrus.Error("执行微信回掉消息类型失败，不存在数据结构体")
		return
	}
	messages := &request.RequestMessageSendMarkdown{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:                 data.FromUser,
			ToParty:                "",
			ToTag:                  "",
			MsgType:                "markdown",
			AgentID:                data.AgentId,
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
	_, err := wecom.WeComApp.Message.SendMarkdown(ctx, messages)
	if err != nil {
		logrus.WithError(err).WithFields(
			logrus.Fields{
				"来自:": data.FromUser,
				"应用":  data.AgentId,
			},
		).Errorf("发送Markdown消息失败")
		return
	}
}
