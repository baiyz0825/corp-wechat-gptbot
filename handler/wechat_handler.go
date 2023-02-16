package handler

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"person-bot/utils"
)

type WxBusStrategy interface {
	Execute(ctx context.Context) error
}

type WxHandler struct {
	WxStrategy map[string]WxBusStrategy
}

func (wx *WxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 校验微信信息
		content := utils.CheckUrlFromWeChat(*r.URL)
		if content != nil {
			// 响应原始数据
			_, err := fmt.Fprintf(w, string(content))
			if err != nil {
				log.Error("写入http请求失败:%w", err)
				return
			}
		}
	} else {
		// 执行正确的业务逻辑
	}
}

// chooseWxStrategy 分发strategy
func (wx *WxHandler) chooseWxStrategy(param string) WxBusStrategy {
	strategy, ok := wx.WxStrategy[param]
	if !ok {
		log.Info("未找到响应的微信命令下的业务逻辑")
	}
	return strategy
}

// registerWxStrategy 注册策略
func (wx *WxHandler) registerWxStrategy(param string, newStrategy WxBusStrategy) {
	if wx.WxStrategy == nil {
		wx.WxStrategy = make(map[string]WxBusStrategy, 20)
	}
	wx.WxStrategy[param] = newStrategy
}

// ParsePostWxMsg 解析请求获取指令数据
func (wx *WxHandler) ParsePostWxMsg(w http.ResponseWriter, r *http.Request) {
	messageEncrypt, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("读取响应数据失败，当前URL是：%v", r.URL)
		w.WriteHeader(http.StatusInternalServerError)
	}
	// 解码
	// 解析与检查url && post body
	realMsgBygte := utils.CheckAndParseBody(*r.URL, string(messageEncrypt))
	// 进行 xml解码
	xml.Unmarshal(realMsg)
	command := utils.ParseCommandFromStr(string(realMsg))
	// 分配处理策略
	ctx := r.Context()
	busContext := context.WithValue(ctx, "msg", realMsg)
	err = wx.chooseWxStrategy(command).Execute(busContext)
	if err != nil {
		log.Error("业务处理失败：%w", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
