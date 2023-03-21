package routers

import (
	"github.com/gin-gonic/gin"
)

func LoadRouters(r *gin.Engine) {

	// 注册微信路由
	RegistryWXRouter(r)
	TestRouter(r)
}
