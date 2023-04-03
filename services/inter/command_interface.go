package inter

import (
	"github.com/baiyz0825/corp-webot/to"
)

type WxTextCommand interface {
	Exec(userData to.MsgContent) bool
}
