package wx_crop

import (
	"encoding/xml"
	"io"
	"net/http"

	"corp-webot/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// VerifyCallBack 回调验证
func VerifyCallBack(r *gin.Context) {
	content := utils.CheckUrlFromWeChat(*r.Request.URL)
	if content == nil {
		logrus.Error("企业微信回调检测失败")
	}
	// 回调检查成功
	logrus.Info("企业微信回调检测成功")
	r.String(http.StatusOK, "%s", content)
}

func RealMsgCallBack(r *gin.Context) {
	buffer, err := io.ReadAll(r.Request.Body)
	if err != nil {

	}
	parsedBytes := utils.CheckAndParseBody(*r.Request.URL, string(buffer))
	// xml解析
	xml.Unmarshal(parsedBytes, obj)
}
