package routers

import (
	"corp-webot/controller/wxcallback"
	"github.com/gin-gonic/gin"
)

func RegistryWXRouter(r *gin.Engine) {
	wxApi := r.Group("/wx")
	{
		wxApi.GET("", wxcallback.VerifyCallBack)
		wxApi.POST("", wxcallback.RealMsgCallBack)
	}
}
