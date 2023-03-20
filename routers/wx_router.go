package routers

import (
	"github.com/baiyz0825/corp-webot/controller"
	"github.com/gin-gonic/gin"
)

func RegistryWXRouter(r *gin.Engine) {
	gptApi := r.Group("/gpt")
	{
		gptApi.GET("", controller.VerifyCallBack)
		gptApi.POST("", controller.ChatWithGPT)
	}
}
