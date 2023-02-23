package http

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"corp-webot/config"
	"github.com/sirupsen/logrus"
)

var HttpClient http.Client

func LoadHttpClientConf() {
	HttpClient = http.Client{
		Timeout: time.Second * 10,
	}
	// 检查是否配置代理
	proxy := config.GetSystemConf().Proxy
	parseUrl, err := url.Parse(proxy)
	if err != nil {
		logrus.Error("代理Url设置错误，本次将不使用代理，请检查代理设置：%w", err)
		HttpClient.Transport = &http.Transport{
			// 跳过证书验证
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:        5000,
			MaxIdleConnsPerHost: 200,
		}
	} else {
		HttpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(parseUrl),
			// 跳过证书验证
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:        5000,
			MaxIdleConnsPerHost: 200,
		}
	}
}
