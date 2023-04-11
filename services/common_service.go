package services

import (
	"encoding/xml"
	"strconv"
	"strings"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
	xcache "github.com/baiyz0825/corp-webot/cache"
	"github.com/baiyz0825/corp-webot/services/command"
	"github.com/baiyz0825/corp-webot/services/event"
	"github.com/baiyz0825/corp-webot/to"
	"github.com/baiyz0825/corp-webot/utils/xlog"
	"github.com/baiyz0825/corp-webot/utils/xstring"
	"github.com/baiyz0825/corp-webot/xconst"
)

// DoTextMsg
// @Description: Text消息逻辑
// @param cmd
func DoTextMsg(userData to.MsgContent) {
	ok := false
	if strings.HasPrefix(userData.Content, xconst.COMMAN_GPT_DELETE_CONTEXT) {
		ok = command.NewContextCommand().Exec(userData)
	} else if strings.HasPrefix(userData.Content, xconst.COMMAND_HELP) {
		ok = command.NewHelpCommandCommand().Exec(userData)
	} else if strings.HasPrefix(userData.Content, xconst.COMMAN_GPT_IMAGE) {
		ok = command.NewGPTImageCommand().Exec(userData)
	} else if strings.HasPrefix(userData.Content, xconst.COMMAN_GPT_PROMPT_SET) {
		ok = command.NewGPTPromptCommand().Exec(userData)
	} else if strings.HasPrefix(userData.Content, xconst.COMMAN_GPT_EXPORT) {
		ok = command.NewExportHistoryCommand().Exec(userData)
	} else {
		ok = command.NewGPTCommand().Exec(userData)
	}
	if !ok {
		xlog.Log.WithField("data:", userData).Error("执行指令失败！")
	}
	xlog.Log.WithField("data:", userData).Debug("执行指令成功！")
	return
}

// DoEventMsg
// @Description: 事件消息逻辑
// @param contentStr
func DoEventMsg(eventData []byte) {
	ok := false
	simpleEvent := &to.SimpleEvent{}
	if err := xml.Unmarshal(eventData, &simpleEvent); err != nil {
		xlog.Log.WithError(err).Error("反序列化用户数据错误")
		return
	}
	switch simpleEvent.Event {
	// 点击事件
	case models.CALLBACK_EVENT_CLICK:
		ok = event.NewClickEventServiceImpl().Exec(eventData)
	default:
		{
			// 不支持返回
			xlog.Log.WithField("data:", simpleEvent).Info("不支持该事件！")
			return
		}
	}
	// 执行失败日志
	if !ok {
		xlog.Log.WithField("data:", simpleEvent).Error("执行指令失败！")
	}
	return
}

// CheckCacheUserEchoReq 检查缓存中是否存在数据 用户多少秒允许发送一次请求
func CheckCacheUserEchoReq(data to.MsgContent) bool {
	hashInt64 := xstring.HashDataConcurrently([]byte(data.Content))
	cacheKey := data.FromUsername + ":" + strconv.FormatInt(hashInt64, 10)
	_, find := xcache.GetCacheDb().Get(cacheKey)
	if find {
		return true
	}
	xcache.SetDataToCache(cacheKey, "", time.Second*4)
	return false
}
