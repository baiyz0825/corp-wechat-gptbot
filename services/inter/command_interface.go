package inter

import (
	"github.com/baiyz0825/corp-webot/to"
)

// CropWxTextCommand
// @Description: 企业微信通用文本指令
type CropWxTextCommand interface {
	Exec(userData to.MsgContent) bool
}
