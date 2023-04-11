package command

import (
	"github.com/baiyz0825/corp-webot/services/impl"
	"github.com/baiyz0825/corp-webot/to"
	"github.com/baiyz0825/corp-webot/xconst"
)

type HelpCommand struct {
}

func NewHelpCommandCommand() *HelpCommand {
	return &HelpCommand{}
}

// Exec
// @Description: 不支持用户指令
// @receiver n
// @param userData
// @return bool
func (n HelpCommand) Exec(userData to.MsgContent) bool {
	impl.SendToWxByText(userData, xconst.GetDefaultNoticeMenu())
	return true
}
