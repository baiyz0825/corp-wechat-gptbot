package handler

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
	"person-bot/config"
	constx "person-bot/const"
	"person-bot/utils"
)

// WeChatConfirmHandler 处理微信回调验证请求
func WeChatConfirmHandler(w http.ResponseWriter, r *http.Request) {
	// 使用url解码
	decode, err := url.Parse(r.URL.String())
	if err != nil {
		log.Error(constx.VerifyS)
	}
	// 获取url请求参数
	param := decode.Query()
	signature := param.Get("msg_signature")
	timestamp := param.Get("timestamp")
	nonce := param.Get("nonce")
	msgEncrypt := param.Get("echostr")
	if !utils.CheckSign(config.GetWechatConf().WeApiRCallToken, signature, timestamp, nonce, msgEncrypt) {
		// 验证失败
		log.Error("企业微信签名验证失败")
		return
	}
	// 解密消息内容
	decrypt, err := utils.MessageDecrypt(msgEncrypt, config.GetWechatConf().WeApiEncodingKey)
	if err != nil {
		// 解密失败
		log.Error("解密回调消息内容失败：%w", err)
		return
	}
	_, _, content, receiverId, err := utils.ParseContent(decrypt)
	if err != nil {
		log.Error("解码回调消息内容失败：%w", err)
		return
	}
	if !strings.EqualFold(string(receiverId), config.GetWechatConf().Corpid) {
		// 不是发给此企业
		log.Error("回调企业错误")
		return
	}
	_, err = io.WriteString(w, string(content))
	if err != nil {
		// 响应错误
		log.Error("写入回调响应错误")
		return
	}
	log.Info("企业微信签名验证成功")
}
