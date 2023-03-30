package wecom

import (
	"context"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/request"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/message/response"
	"github.com/baiyz0825/corp-webot/config"
	"github.com/baiyz0825/corp-webot/to"
	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
	"github.com/sirupsen/logrus"
)

var wxCrypt *wxbizmsgcrypt.WXBizMsgCrypt

var WeComApp *work.Work

func init() {
	LoadWeComAppConf()
	LoadWxUtils()
}

func LoadWeComAppConf() {
	logrus.Info("初始化企业微信助手......")
	app, err := work.NewWork(&work.UserConfig{
		CorpID:  config.GetWechatConf().Corpid,     // 企业微信的app id，所有企业微信共用一个。
		AgentID: config.GetWechatConf().AgentId,    // 内部应用的app id
		Secret:  config.GetWechatConf().CorpSecret, // 内部应用的app secret
		OAuth: work.OAuth{
			Callback: config.GetSystemConf().CallBackUrl, //
			Scopes:   nil,
		},
		HttpDebug: true,
		// 可选，不传默认走程序内存
		// Cache: kernel.NewRedisClient(&kernel.RedisOptions{
		// 	Addr:     "127.0.0.1:6379",
		// 	Password: "",
		// 	DB:       0,
		// }),
	})
	if err != nil {
		logrus.WithError(err).Error("初始化企业微信助手失败！")
		panic(err)
	}
	WeComApp = app
}

func LoadWxUtils() {
	logrus.Info("初始化微信工具包......")
	wxCrypt = wxbizmsgcrypt.NewWXBizMsgCrypt(config.GetWechatConf().WeApiRCallToken, config.GetWechatConf().WeApiEncodingKey, config.GetWechatConf().Corpid, wxbizmsgcrypt.XmlType)
}

// GetReVerifyCallBack 从微信回调解析请求数据
func GetReVerifyCallBack(q to.CallBackParams) []byte {
	msg, cryptErr := wxCrypt.VerifyURL(q.MsgSignature, q.TimeStamp, q.Nonce, q.Echostr)
	if cryptErr != nil {
		logrus.Errorf("验证Url出错（回调消息解密错误）：%v", cryptErr)
		return []byte("")
	}
	logrus.Info("解析的回调字符为：", string(msg))
	return msg
}

// DeCryptMsg 解密消息
func DeCryptMsg(cryptMsg []byte, msgSignature, timeStamp, nonce string) []byte {
	msg, cryptErr := wxCrypt.DecryptMsg(msgSignature, timeStamp, nonce, cryptMsg)
	if cryptErr != nil {
		logrus.Errorf("回调消息解密错误：%v", cryptErr)
		return nil
	}
	return msg
}

// CryptMessage 加密消息
func CryptMessage(respData, reqTimestamp, reqNonce string) string {
	encryptMsg, cryptErr := wxCrypt.EncryptMsg(respData, reqTimestamp, reqNonce)
	if cryptErr != nil {
		logrus.Errorf("消息加密错误：%v", cryptErr)
		return ""
	}
	return string(encryptMsg)
}

// SendMarkdownToUSer 发送Markdown消息
func SendMarkdownToUSer(data to.MsgContent, respMsg string) *response.ResponseMessageSend {
	// 封装微信消息体
	messages := &request.RequestMessageSendMarkdown{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:                 data.FromUsername,
			ToParty:                "",
			ToTag:                  "",
			MsgType:                "markdown",
			AgentID:                config.GetWechatConf().AgentId,
			EnableDuplicateCheck:   1,
			DuplicateCheckInterval: 1800,
		},
		Markdown: &request.RequestMarkdown{
			Content: respMsg,
		},
	}
	// 发送微信消息
	resp, err := WeComApp.Message.SendMarkdown(context.Background(), messages)
	if err != nil {
		logrus.Errorf("创建微信发送消息内容失败：%v", err)
		return nil
	}
	return resp
}

// SendMarkdownToUSer 发送Markdown消息
func SendTextToUSer(data to.MsgContent, respMsg string) *response.ResponseMessageSend {
	// 封装微信消息体
	messages := &request.RequestMessageSendText{
		RequestMessageSend: request.RequestMessageSend{
			ToUser:                 data.FromUsername,
			ToParty:                "",
			ToTag:                  "",
			MsgType:                "text",
			AgentID:                config.GetWechatConf().AgentId,
			Safe:                   0,
			EnableIDTrans:          0,
			EnableDuplicateCheck:   1,
			DuplicateCheckInterval: 1800,
		},
		Text: &request.RequestText{
			Content: respMsg,
		},
	}
	// 发送微信消息
	resp, err := WeComApp.Message.SendText(context.Background(), messages)
	if err != nil {
		logrus.Errorf("创建微信发送消息内容失败：%v", err)
		return nil
	}
	return resp
}
