package inter

import (
	"github.com/baiyz0825/corp-webot/to"
)

// SystemCmd
// @Description: 企业微信通用文本指令
type SystemCmd interface {
	Exec(userData to.MsgContent) bool
}

// CropEvent
// @Description: 事件处理
type CropEvent interface {
	Exec(eventData []byte, userData to.SimpleEvent) bool
}
