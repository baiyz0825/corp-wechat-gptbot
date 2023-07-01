package xhttp

import (
	"crypto/tls"
	"fmt"
	"net"
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
		if CheckServer(parseUrl.Host) && err == nil {
			xlog.Log.Info("代理Url获取成功，本次将使用代理")
			HttpClient.Transport = &http.Transport{
				Timeout:         time.Second * 120,
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				Proxy:           http.ProxyURL(parseUrl),
			}
			return
		}
	} 
	// 未设置代理或者不可用
	xlog.Log.Infof("客户端Http代理未设置设置，本次将不使用代理")
	HttpClient.Transport = &http.Transport{
		Timeout:         time.Second * 120,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return
}

func CheckServer(strUrl string) bool {
	timeout := 5 * time.Second
	_, err := net.DialTimeout("tcp", strUrl, timeout)
	if err != nil {
		fmt.Println("无法访问代理, error: ", err)
		return false
	}
	xlog.Log.Infof("代理测试成功！")
	return true
}
