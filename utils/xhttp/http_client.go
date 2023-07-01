package xhttp

import (
	"crypto/tls"
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
			HttpClient.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				Proxy:           http.ProxyURL(parseUrl),
			}
			// 检查代理是否可用
			CheckProxy()
			return
		}
	}
	// 未设置代理
	xlog.Log.Infof("未设置代理启动!")
	HttpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return
}

// CheckProxy 测试代理
func CheckProxy() {
	// 发起GET请求测试代理
	resp, err := HttpClient.Get("https://www.google.com/")
	if err != nil {
		// 处理请求错误
		xlog.Log.Errorf("测试代理请求错误:%s", err)
		xlog.Log.Fatalf("代理检测失败！请重新配置！")
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		xlog.Log.Infof("代理可用，设置代理启动!")
	}
	return
}
