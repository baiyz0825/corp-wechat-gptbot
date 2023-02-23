package routers

import (
	"corp-webot/controller/chatgpt"
	"github.com/gin-gonic/gin"
)

func RegistryWXRouter(r *gin.Engine) {
	wxApi := r.Group("/wx")
	{
		wxApi.GET("", chatgpt.VerifyCallBack)
		wxApi.POST("", chatgpt.RealMsgCallBack)
	}
}
