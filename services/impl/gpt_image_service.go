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
	// 获取消息命令信息
	userData.Content = strings.TrimPrefix(userData.Content, g.Command)
	// 调用api生成图像
	url := openaiutils.SendReqAndGetImageResp(userData.Content)
	if len(url) == 0 {
		return false
	}
	// 下载图像 -> 微信临时素材
	bytes, fileType, err := xhttp.DownloadImageGetBytes(url, xhttp.HttpClient)
	if err != nil {
		xlog.Log.WithError(err).WithField("请求用户是：", userData.ToUsername).Error("下载openai图片失败")
		return false
	}
	// 响应数据回用户端
	resp := wecom.SendImageToUser(bytes, fileType, userData.ToUsername)
	if resp == nil {
		return false
	}
	return true
}
