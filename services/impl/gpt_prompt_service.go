package impl

import (
	"strings"

	"github.com/baiyz0825/corp-webot/dao"
	"github.com/baiyz0825/corp-webot/to"
	"github.com/baiyz0825/corp-webot/utils/xlog"
	"github.com/baiyz0825/corp-webot/xconst"
)

type GPTPromptCommand struct {
	Command string
}

func NewGPTPromptCommand() *GPTPromptCommand {
	return &GPTPromptCommand{Command: xconst.COMMAN_GPT_PROMPT_SET}
}

func (g GPTPromptCommand) Exec(userData to.MsgContent) bool {
	userData.Content = strings.TrimPrefix(userData.Content, g.Command)
	// update db
	err := dao.UpdateUser(userData.Content, userData.FromUsername)
	if err != nil {
		xlog.Log.WithError(err).WithField("用户:", userData.FromUsername).Error("设置用户自定义提示词失败")
		SendToWxByText(userData, xconst.AI_DEFAULT_MSG)
		return false
	}
	SendToWxByText(userData, xconst.AI_KNOWN_YOUR_ASK)
	return true
}
