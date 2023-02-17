package utils

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"time"

	"corp-webot/config"
	"github.com/sirupsen/logrus"
)

var httpClient *http.Client

func LoadHttpClientConf() {
	httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	// 检查是否配置代理
	proxy := config.GetSystemConf().Proxy
	parseUrl, err := url.Parse(proxy)
	if err != nil {
		logrus.Error("代理Url设置错误，本次将不使用代理，请检查代理设置：%v", err)
		httpClient.Transport = &http.Transport{
			// 跳过证书验证
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:        5000,
			MaxIdleConnsPerHost: 200,
		}
	} else {
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(parseUrl),
			// 跳过证书验证
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:        5000,
			MaxIdleConnsPerHost: 200,
		}
	}
}

func Get(urlStr string, headers map[string]string) ([]byte, error) {
	request, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	// 检查header
	if headers != nil {
		for key, value := range headers {
			request.Header.Set(key, value)
		}
	}
	// 发送请求
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	// 关闭响应
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error("http client close body err : %v", err)
		}
	}(response.Body)
	// 获取数据
	return nil, err
}

func Post() {

}

func Gets() {

}

func Posts() {

}
