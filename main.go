package main

import (
	"github.com/baiyz0825/corp-webot/config"
	"github.com/baiyz0825/corp-webot/middleware"
	"github.com/baiyz0825/corp-webot/routers"
	"github.com/gin-gonic/gin"
)

func loadGin() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	// 使用中间件
	r.Use(middleware.LoggerToFile())
	// 注册路由
	routers.LoadRouters(r)
	return r
}

func main() {
	r := loadGin()
	// 启动gin
	r.Run(":" + config.GetSystemConf().Port)
}
