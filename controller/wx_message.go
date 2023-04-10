package controller

import (
	"encoding/xml"
	"net/http"

	"github.com/baiyz0825/corp-webot/dao"
	"github.com/baiyz0825/corp-webot/services/impl"
	"github.com/baiyz0825/corp-webot/to"
	"github.com/baiyz0825/corp-webot/utils/wecom"
	"github.com/baiyz0825/corp-webot/utils/xlog"
	"github.com/gin-gonic/gin"
)

// VerifyCallBack 回调验证
func VerifyCallBack(c *gin.Context) {
	var q to.CallBackParams
	if err := c.Bind(&q); err != nil {
		xlog.Log.Errorf("绑定回调Query错误：%v", err)
	}
	msg := wecom.GetReVerifyCallBack(q)
	_, _ = c.Writer.Write(msg)
}

// WxChatCommand 实际处理用户消息
func WxChatCommand(c *gin.Context) {
	var dataStuc to.CallBackData
	if err := c.ShouldBindQuery(&dataStuc); err != nil {
		xlog.Log.Errorf("绑定回调Query错误：%v", err)
	}
	// 解析请求体
	raw, err := c.GetRawData()
	if err != nil {
		xlog.Log.WithError(err).Error("解析微信回调参数失败")
		return
	}
	userData := to.MsgContent{}
	userDataDecrypt := wecom.DeCryptMsg(raw, dataStuc.MsgSignature, dataStuc.TimeStamp, dataStuc.Nonce)
	// 解密失败返回空
	if userDataDecrypt == nil {
		xlog.Log.WithField("用户数据：", userData).Error("解密失败")
	}
	// 提前向微信返回成功接受，防止微信多次回调
	c.JSON(http.StatusOK, "")
	// 异步处理用户请求
	go func() {
		err = xml.Unmarshal(userDataDecrypt, &userData)
		// 检测缓存
		if impl.CheckCacheUserEchoReq(userData) {
			return
		}
		if err != nil {
			xlog.Log.WithError(err).Error("反序列化用户数据错误")
			return
		}
		if userData.MsgType != "text" {
			// 已经返回了数据类型，这里不能返回数据，会导致错误
			// TODO 菜单逻辑
			// c.String(http.StatusBadRequest, "不支持非text类型处理")
			return
		}
		// 检查用户是否存在，不存在创建
		if !dao.CheckUserAndCreate(userData.FromUsername) {
			xlog.Log.WithField("用户信息", userData.FromUsername).Errorf("创建用户失败")
			return
		}
		// 处理数据
		ok := impl.GetCommand(userData.Content).Exec(userData)
		if !ok {
			xlog.Log.WithField("data:", userData).Error("执行指令失败！")
			return
		}
		xlog.Log.WithField("data:", userData).Debug("执行指令成功！")
		return
	}()
}
