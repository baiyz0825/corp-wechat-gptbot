package routers

import (
	"corp-webot/controller/message/wx_crop"
	"github.com/gin-gonic/gin"
)

func RegistryWXRouter(r *gin.Engine) {
	wxApi := r.Group("/wx")
	{
		wxApi.GET("", wx_crop.VerifyCallBack)
		wxApi.POST("", wx_crop.RealMsgCallBack)
	}
}
