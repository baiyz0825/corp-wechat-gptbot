package controller

import (
	"encoding/xml"
	"net/http"

	"github.com/baiyz0825/corp-webot/services"
	"github.com/baiyz0825/corp-webot/to"
	"github.com/baiyz0825/corp-webot/utils/wecom"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// VerifyCallBack 回调验证
func VerifyCallBack(c *gin.Context) {
	var q to.CallBackParams
	if err := c.Bind(&q); err != nil {
		logrus.Errorf("绑定回调Query错误：%v", err)
	}
	msg := wecom.GetReVerifyCallBack(q)
	_, _ = c.Writer.Write(msg)
}

// ChatWithGPT 实际处理用户消息
func ChatWithGPT(c *gin.Context) {
	var dataStuc to.CallBackData
	if err := c.ShouldBindQuery(&dataStuc); err != nil {
		logrus.Errorf("绑定回调Query错误：%v", err)
	}
	// 解析请求体
	raw, err := c.GetRawData()
	if err != nil {
		logrus.WithError(err).Error("解析微信回调参数失败")
		return
	}
	userData := to.MsgContent{}
	userDataDecrypt := wecom.DeCryptMsg(raw, dataStuc.MsgSignature, dataStuc.TimeStamp, dataStuc.Nonce)
	// 解密失败返回空
	if userDataDecrypt == nil {
		logrus.WithField("用户数据：", userData).Error("解密失败")
	}
	// 提前向微信返回成功接受，防止微信多次回调
	c.JSON(http.StatusOK, "")
	// 异步处理用户请求
	go func() {
		err = xml.Unmarshal(userDataDecrypt, &userData)
		// 检测缓存
		if services.CheckCacheUserEchoReq(userData) {
			return
		}
		if err != nil {
			logrus.WithError(err).Error("反序列化用户数据错误")
			return
		}
		if userData.MsgType != "text" {
			c.String(http.StatusBadRequest, "不支持非text类型处理")
		}
		// 处理数据
		ok := services.DoChat(userData)
		if !ok {
			logrus.WithField("data:", userData).Error("发送失败")
			return
		}
		logrus.WithField("data:", userData).Error("发送成功")
		return
	}()
}
