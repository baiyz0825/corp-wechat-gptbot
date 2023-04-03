package xlog

import (
	"io"
	"os"

	"github.com/baiyz0825/corp-webot/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = NewLogger()
}

func NewLogger() *logrus.Logger {
	Logger := logrus.New()
	// 日志级别
	logLev, err := logrus.ParseLevel(config.GetSystemConf().LogConf.LogLevel)
	Logger.SetLevel(logLev)
	if err != nil {
		Logger.SetLevel(logrus.DebugLevel)
	}
	// writer
	switch config.GetSystemConf().LogConf.LogOutPutMode {
	case "console":
		Logger.SetOutput(os.Stdout)
	case "file":
		path := config.GetSystemConf().LogConf.LogOutPutPath + "_%Y%m%d%H%M" + ""
		writer, _ := rotatelogs.New(
			path,
			rotatelogs.WithLinkName(config.GetSystemConf().LogConf.LogOutPutPath),
			rotatelogs.WithRotationSize(config.GetSystemConf().LogConf.LogFileMaxSizeM),
			rotatelogs.WithRotationCount(config.GetSystemConf().LogConf.LogFileRotationCount),
		)
		Logger.SetOutput(writer)
	case "both":
		{
			path := config.GetSystemConf().LogConf.LogOutPutPath + "_%Y%m%d%H%M" + ""
			writer, _ := rotatelogs.New(
				path,
				rotatelogs.WithLinkName(config.GetSystemConf().LogConf.LogOutPutPath),
				rotatelogs.WithRotationSize(config.GetSystemConf().LogConf.LogFileMaxSizeM),
				rotatelogs.WithRotationCount(config.GetSystemConf().LogConf.LogFileRotationCount),
			)
			multiWriter := io.MultiWriter(os.Stdout, writer)
			Logger.SetOutput(multiWriter)
		}
	default:
	}
	// 格式化
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: config.GetSystemConf().LogConf.LogFileDateFmt,
	})
	if config.GetSystemConf().LogConf.LogFormatter == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: config.GetSystemConf().LogConf.LogFileDateFmt,
		})
	}
	return Logger
}
