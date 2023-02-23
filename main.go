package main

import (
	"time"

	"corp-webot/cache"
	"corp-webot/config"
	"corp-webot/middleware"
	"corp-webot/routers"
	"corp-webot/utils/gpt"
	"corp-webot/utils/http"
	"corp-webot/utils/wecom"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
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

// CronProcess 定时任务
func CronProcess() {
	location, err := time.LoadLocation(" Asia/Shanghai")
	if err != nil {
		return
	}
	cronJob := cron.New(cron.WithLocation(location))
	cronJob.AddFunc("0 */2 * * *", func() {
		logrus.Info("正在执行刷新GPT接口AccessToken定时任务")
		gpt.LoadAccessToken()
	})
	// 开始定时任务
	cronJob.Start()
	defer cronJob.Stop()
}

func InitLoadComponent() {
	// 加载依赖工具
	cache.LoadCache()
	http.LoadHttpClientConf()
	wecom.LoadWxBizCryptHelper()
	wecom.LoadWeComAppConf()
	gpt.LoadGptUtils()
}

func main() {
	r := loadGin()
	// load conf
	if err := config.LoadConf(); err != nil {
		logrus.Fatal(err)
	}
	InitLoadComponent()
	// 开始定时任务
	CronProcess()
	// 启动gin
	r.Run(":" + config.GetSystemConf().Port)
}
