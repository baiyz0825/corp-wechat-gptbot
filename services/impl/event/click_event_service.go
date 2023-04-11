package event

import (
	"encoding/xml"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
	cmd "github.com/baiyz0825/corp-webot/services/impl/command"
	"github.com/baiyz0825/corp-webot/to"
	"github.com/baiyz0825/corp-webot/utils/xlog"
	"github.com/baiyz0825/corp-webot/xconst"
)

type ClickEventServiceImpl struct {
	Event string
}

func NewClickEventServiceImpl() *ClickEventServiceImpl {
	return &ClickEventServiceImpl{Event: models.CALLBACK_EVENT_CLICK}
}

func (c ClickEventServiceImpl) Exec(eventData []byte, userData to.MsgContent) bool {
	// 序列化数据
	clickEvent := &models.EventClick{}
	err := xml.Unmarshal(eventData, clickEvent)
	// 解密失败返回空
	if err != nil {
		xlog.Log.WithField("用户事件：", c.Event).Error("序列化事件失败")
		return false
	}
	// 分发点击事件功能
	switch clickEvent.EventKey {
	case xconst.COMMAN_GPT_DELETE_CONTEXT:
		return cmd.NewContextCommand().Exec(userData)
	case xconst.COMMAN_GPT_EXPORT:
		return cmd.NewExportHistoryCommand().Exec(userData)
	case xconst.COMMAND_HELP:
		return cmd.NewHelpCommandCommand().Exec(userData)
	default:
		// 默认分发成功
		return true
	}
}
