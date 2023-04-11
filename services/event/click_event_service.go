package event

import (
	"encoding/xml"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
	"github.com/baiyz0825/corp-webot/services/command"
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

func (c ClickEventServiceImpl) Exec(eventData []byte) bool {
	// 序列化数据
	clickEvent := &models.EventClick{}
	err := xml.Unmarshal(eventData, clickEvent)
	// 解密失败返回空
	if err != nil {
		xlog.Log.WithField("用户事件：", c.Event).Error("序列化事件失败")
		return false
	}
	userData := to.MsgContent{
		ToUsername:   clickEvent.ToUserName,
		FromUsername: clickEvent.FromUserName,
		MsgType:      clickEvent.MsgType,
		Agentid:      clickEvent.AgentID,
	}
	// 分发点击事件功能
	switch clickEvent.EventKey {
	case xconst.COMMAN_GPT_DELETE_CONTEXT:
		return command.NewContextCommand().Exec(userData)
	case xconst.COMMAN_GPT_EXPORT:
		return command.NewExportHistoryCommand().Exec(userData)
	case xconst.COMMAND_HELP:
		return command.NewHelpCommandCommand().Exec(userData)
	default:
		// 默认分发成功
		return true
	}
}
