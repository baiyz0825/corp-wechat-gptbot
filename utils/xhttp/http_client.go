package xhttp

import (
	"net/http"
	"net/url"
	"time"

	"github.com/baiyz0825/corp-webot/config"
	"github.com/sirupsen/logrus"
)

var HttpClient http.Client

func init() {
	logrus.Info("初始化HTTP客户端......")
	HttpClient = http.Client{
		Timeout: time.Second * 60,
	}
	// 检查是否配置代理
	proxy := config.GetSystemConf().Proxy
	if len(proxy) > 0 {
		parseUrl, err := url.Parse(proxy)
		if err != nil {
			logrus.Error("代理Url获取解析失败，本次将不使用代理")
		}
		logrus.Infof("设置代理中...")
		HttpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(parseUrl),
		}
		return
	}
	logrus.Infof("客户端Http代理未设置设置，本次将不使用代理")
	return
}
