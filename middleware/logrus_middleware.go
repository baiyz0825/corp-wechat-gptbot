package middleware

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// LoggerToFile 日志记录到文件
func LoggerToFile() gin.HandlerFunc {
	logFilePath := ""
	logFileName := ""
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	// 写入文件
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	// 创建多个writer，同时写入控制台文件
	writes := []io.Writer{
		file,
		os.Stdout,
	}
	multiWriters := io.MultiWriter(writes...)
	// 实例化
	logger = logrus.New()
	// 设置输出
	logger.SetOutput(multiWriters)
	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	// 设置日志格式TXT
	logger.SetFormatter(&logrus.TextFormatter{})
	// 设置日志格式JSON
	// logger.SetFormatter(&logrus.JSONFormatter{
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// })
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 日志格式
		// logger.Infof("| %3d | %13v | %15s | %s | %s |",
		// 	statusCode,
		// 	latencyTime,
		// 	clientIP,
		// 	reqMethod,
		// 	reqUri,
		// )
		// 日志格式
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
		}).Info()
	}
}

// LoggerToMongo 日志记录到 MongoDB
func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

// LoggerToES 日志记录到 ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

// LoggerToMQ 日志记录到 MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
