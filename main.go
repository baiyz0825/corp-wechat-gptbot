package main

import (
	"github.com/baiyz0825/corp-webot/config"
	"github.com/baiyz0825/corp-webot/dao"
	"github.com/baiyz0825/corp-webot/middleware"
	"github.com/baiyz0825/corp-webot/routers"
	"github.com/baiyz0825/corp-webot/utils/xstring"
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
	dao.LoadDatabase()
	return r
}

func main() {
	r := loadGin()
	// logo
	xstring.GenLogoAscii("GPT-BOT", "green")
	// 启动gin
	_ = r.Run(":" + config.GetSystemConf().Port)
	// 关闭后置操作
	dao.CloseDb()
}
