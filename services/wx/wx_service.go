package wx

import (
	"strings"

	"github.com/baiyz0825/corp-webot/to"
	"github.com/baiyz0825/corp-webot/utils/wecom"
	"github.com/sirupsen/logrus"
)

// SendToWxByText 使用大文本发送
func SendToWxByText(userData to.MsgContent, msg string) bool {
	// 按行分割代码
	lines := strings.Split(msg, "\n")
	// 临时消息缓冲
	var content string
	for i := 0; i < len(lines); i++ {
		// >2000 发送前一部分，清空重来
		if len(content)+len(lines[i]) > 2000 {
			resp := wecom.SendTextToUSer(userData.FromUsername, content)
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
		resp := wecom.SendTextToUSer(userData.FromUsername, content)
		if resp.ResponseWork.ErrCode != 0 {
			logrus.WithField("resp:", resp).Error("企业微信助手发送分片失败")
		}
	}
	return true
}

// SendToWxByMarkdown 使用markdown发送
func SendToWxByMarkdown(userData to.MsgContent, msg string) bool {
	// TODO: 考虑是否分片长消息分割
	resp := wecom.SendMarkdownToUSer(userData.FromUsername, msg)
	if resp.ResponseWork.ErrCode != 0 {
		logrus.WithField("resp:", resp).Error("企业微信助手发送失败")
		return false
	}
	return true
}
