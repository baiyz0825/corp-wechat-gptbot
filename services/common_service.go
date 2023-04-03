package services

import (
	"strings"
	"time"

	xcache "github.com/baiyz0825/corp-webot/cache"
	"github.com/baiyz0825/corp-webot/services/impl"
	"github.com/baiyz0825/corp-webot/services/inter"
	"github.com/baiyz0825/corp-webot/to"
	"github.com/baiyz0825/corp-webot/utils/wecom"
	"github.com/baiyz0825/corp-webot/xconst"
	"github.com/sirupsen/logrus"
)

type Command struct {
	cmd inter.WxTextCommand
}

func (n Command) Exec(userData to.MsgContent) bool {
	return n.cmd.Exec(userData)
}

// GetCommand
// @Description: 获取指令服务
// @param cmd
func GetCommand(contentStr string) *Command {
	command := &Command{}
	if strings.HasPrefix(contentStr, xconst.COMMAND_GPT) {
		command.cmd = impl.NewGPTCommand()
	} else if strings.HasPrefix(contentStr, xconst.COMMAN_GPT_DELETE_CONTEXT) {
		command.cmd = impl.NewContextCommand()
	} else {
		command.cmd = impl.NewNotSupportCommand()
	}
	return command
}

// SendToWxByMarkdown 使用markdown发送
func SendToWxByMarkdown(userData to.MsgContent, msg string) bool {
	// TODO: 考虑是否分片长消息分割
	resp := wecom.SendMarkdownToUSer(userData, msg)
	if resp.ResponseWork.ErrCode != 0 {
		logrus.WithField("resp:", resp).Error("企业微信助手发送失败")
		return false
	}
	return true
}

// SendToWxByText 使用大文本发送
func SendToWxByText(userData to.MsgContent, msg string) bool {
	// 按行分割代码
	lines := strings.Split(msg, "\n")
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
