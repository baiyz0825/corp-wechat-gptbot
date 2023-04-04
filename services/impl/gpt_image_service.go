package impl

import (
	"strings"

	"github.com/baiyz0825/corp-webot/to"
	"github.com/baiyz0825/corp-webot/utils/openaiutils"
	"github.com/baiyz0825/corp-webot/utils/wecom"
	"github.com/baiyz0825/corp-webot/utils/xhttp"
	"github.com/baiyz0825/corp-webot/utils/xlog"
	"github.com/baiyz0825/corp-webot/xconst"
)

type GPTImageCommand struct {
	Command string
}

func NewGPTImageCommand() *GPTImageCommand {
	return &GPTImageCommand{Command: xconst.COMMAN_GPT_IMAGE}
}

// Exec
// @Description: 图片生成逻辑
// @receiver g
// @param userData
// @return bool
func (g GPTImageCommand) Exec(userData to.MsgContent) bool {
	// 调用api生成图像
	xlog.Log.WithField("用户:", userData.FromUsername).WithField("请求数据是:", userData.Content).Debug("开始请求openai接口生成图像....")
	url := openaiutils.SendReqAndGetImageResp(strings.TrimPrefix(userData.Content, g.Command))
	if len(url) == 0 {
		wecom.SendTextToUSer(userData.FromUsername, xconst.AI_API_ERROR_MSG)
		return false
	}
	xlog.Log.WithField("用户:", userData.FromUsername).WithField("请求数据是:", url).Debug("开始下载图像....")
	// 下载图像 -> 微信临时素材
	bytes, fileType, err := xhttp.DownloadImageGetBytes(url, xhttp.HttpClient)
	if err != nil {
		xlog.Log.WithError(err).WithField("请求用户是：", userData.FromUsername).Error("下载openai图片失败")
		return false
	}
	xlog.Log.WithField("用户:", userData.FromUsername).WithField("请求数据是:", url).Debug("开始上传微信素材，并发送给用户....")
	// 响应数据回用户端
	resp := wecom.SendImageToUser(bytes, fileType, userData.FromUsername)
	if resp == nil {
		return false
	}
	return true
}
