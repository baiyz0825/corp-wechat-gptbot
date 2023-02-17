package wx_crop

import (
	"encoding/xml"
	"net/http"

	"corp-webot/model"
	"corp-webot/utils"
	"corp-webot/xconst"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type QueryParams struct {
	MsgSignature string `form:"msg_signature"`
	TimeStamp    string `form:"timestamp"`
	Nonce        string `form:"nonce"`
	EchoStr      string `form:"echostr"`
}

// VerifyCallBack 回调验证
func VerifyCallBack(r *gin.Context) {
	var q QueryParams
	if err := r.Bind(&q); err != nil {
		logrus.Errorf("绑定回调Query错误：%v", err)
	}
	echoStr, cryptErr := utils.WxBizMsgCryptHelper.VerifyURL(q.MsgSignature, q.TimeStamp, q.Nonce, q.EchoStr)
	if cryptErr != nil {
		logrus.Errorf("验证Url出错：%v", cryptErr)
	}
	logrus.Info("解析的回调字符为：", string(echoStr))
	r.Writer.Write(echoStr)
}

func RealMsgCallBack(r *gin.Context) {
	q := QueryParams{
		MsgSignature: r.Query("msg_signature"),
		TimeStamp:    r.Query("timestamp"),
		Nonce:        r.Query("nonce"),
		EchoStr:      "",
	}
	// 验证url请求
	rawData, err := r.GetRawData()
	if err != nil {
		logrus.Errorf("解析请求数据错误：%v", err)
		return
	}
	msgStr, cryptError := utils.WxBizMsgCryptHelper.DecryptMsg(q.MsgSignature, q.TimeStamp, q.Nonce, rawData)
	if cryptError != nil {
		logrus.Errorf("校验获取回调数据失败：%v", err)
		return
	}
	// 正确解析
	var textMessage model.RecTextMessage
	err = xml.Unmarshal(msgStr, &textMessage)
	if err != nil {
		logrus.Errorf("xml映射消息对象失败：%v", err)
		r.Status(http.StatusInternalServerError)
		return
	}
	// 响应数据
	r.Status(http.StatusOK)
	// 结束当前处理，调用goroutine处理
	r.Done()
	// 分发策略处理
	go func(msg model.RecTextMessage) {
		// 获取用户传递标记
		command := string([]rune(msg.Content)[0:4])
		_, ok := xconst.TextCommandMap[command]
		if !ok {
			// 给用户回复消息，指令列表
		}
		// 分发策略处理

	}(textMessage)
}
