package main

import (
	"corp-webot/config"
	"corp-webot/middleware"
	"corp-webot/routers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	// load conf
	if err := config.LoadConf(); err != nil {
		logrus.Fatal(err)
	}
	// Listen and Server in 0.0.0.0:8080
	r.Run(":" + config.GetSystemConf().Port)
}
