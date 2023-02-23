package services

import (
	"time"

	"corp-webot/cache"
	gpt2 "corp-webot/utils/gpt"
	"corp-webot/utils/http"
	string2 "corp-webot/utils/string"
	"corp-webot/utils/wecom"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type GPTChatCommand struct {
}

func (gpt *GPTChatCommand) ExecCommand(data *CommandData, ctx context.Context) {
	// 获取用户请求中的消息
	if data != nil {
		logrus.Error("执行微信回掉文本命令失败，不存在命令结构体")
		return
	}

	// 查找用户对应的gpt
	var gptHelper gpt2.V1GPTHelper
	if gptHelper = cache.GetDataFromCache(data.Cmd + data.FromUser + data.CorpID).(gpt2.V1GPTHelper); &gptHelper == nil {
		logrus.Info("当前用户不存在GptHelper，准备重建")
		gptHelper = gpt2.NewGptHelperV1(gpt2.Aksk, data.Cmd+data.FromUser+data.CorpID)
		// 放入缓存
		cache.SetDataToCache(data.Cmd+data.FromUser+data.CorpID, &gptHelper, time.Minute*1)
	}
	logrus.WithField("当前TraceId", gptHelper.GetTraceIDFromHelper()).Debugf("解析用户" + data.FromUser + "发送请求成功，准备请求GPT")
	// 发送到ChatGPT进行解析
	respMsg, err := gptHelper.SendAndGetMessageToGPTV1(http.HttpClient, data.Msg)
	if err != nil {
		logrus.WithError(err).WithField("当前TraceId", gptHelper.GetTraceIDFromHelper()).Errorf("向GPT发送当前用户" + data.FromUser + "请求消息失败！")
		return
	}
	logrus.WithField("当前TraceId", gptHelper.GetTraceIDFromHelper()).Debugf("用户" + data.FromUser + "已从GPT获取响应")
	// 获取解析结果,转换为markdown文本
	markdownMsg := string2.TransBytesToMarkdownStr(string(respMsg))
	// 封装微信消息体
	messages := &request.RequestMessageSendMarkdown{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:                 "UserID1|UserID2|UserID3",
			ToParty:                "PartyID1|PartyID2",
			ToTag:                  "TagID1 | TagID2",
			MsgType:                "markdown",
			AgentID:                1,
			EnableDuplicateCheck:   0,
			DuplicateCheckInterval: 1800,
		},
		Markdown: &request.RequestMarkdown{
			Content: markdownMsg,
		},
	}
	// 发送微信消息
	_, err = wecom.WeComApp.Message.SendMarkdown(ctx, messages)
	if err != nil {
		logrus.WithError(err).WithField("当前TraceId", gptHelper.GetTraceIDFromHelper()).Errorf("发送Markdown消息失败")
		return
	}
}
