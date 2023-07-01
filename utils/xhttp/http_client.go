package xhttp

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/baiyz0825/corp-webot/config"
	"github.com/baiyz0825/corp-webot/utils/xlog"
)

var HttpClient *http.Client

func init() {
	xlog.Log.Info("初始化HTTP客户端......")
	HttpClient = &http.Client{
		Timeout: time.Second * 60,
	}
	// 检查是否配置代理
	proxy := config.GetSystemConf().Proxy
	if len(proxy) > 0 {
		parseUrl, err := url.Parse(proxy)
		if err == nil {
			xlog.Log.Info("代理Url获取成功，本次将使用代理")
			HttpClient.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				Proxy:           http.ProxyURL(parseUrl),
			}
		}
	}
	// 未设置代理或者不可用
	CheckProxy()
	return
}

// CheckProxy 测试代理
func CheckProxy() {
	// 发起GET请求测试代理
	resp, err := HttpClient.Get("https://www.google.com/")
	if err != nil {
		// 处理请求错误
		xlog.Log.Infof("测试代理请求错误:%s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)
	if resp.StatusCode == http.StatusOK {
		xlog.Log.Infof("代理可用")
		return
	}
	xlog.Log.Infof("代理不可用，状态码:%d", resp.StatusCode)
	xlog.Log.Infof("客户端Http代理未设置或不能使用，***** 本次将不使用代理启动 *****")
	HttpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return
}
