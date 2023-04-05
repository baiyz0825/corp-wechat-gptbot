package impl

import (
	"bytes"
	"fmt"
	"time"

	"github.com/baiyz0825/corp-webot/dao"
	"github.com/baiyz0825/corp-webot/model"
	"github.com/baiyz0825/corp-webot/to"
	"github.com/baiyz0825/corp-webot/utils/doc"
	"github.com/baiyz0825/corp-webot/utils/wecom"
	"github.com/baiyz0825/corp-webot/utils/xlog"
	"github.com/baiyz0825/corp-webot/xconst"
	"github.com/pkg/errors"
)

type ExportHistoryCommand struct {
	Command string
}

func NewExportHistoryCommand() *ExportHistoryCommand {
	return &ExportHistoryCommand{Command: xconst.COMMAN_GPT_EXPORT}
}

func (e ExportHistoryCommand) Exec(userData to.MsgContent) bool {
	// 获取用户名
	username := userData.FromUsername
	// 查询db获取上下文
	context, err := dao.GetLatestUserContext(username)
	if err != nil || context == nil {
		xlog.Log.WithField("用户:", username).WithError(err).Error("获取导出用户历史上下文数据失败！")
		return false
	}
	// 生成导出文件
	pdfBytes, err := SearchHistoryAndGenPDF(context, err, username)
	if err != nil {
		xlog.Log.WithField("用户:", username).WithError(err).Error("获取导出用户历史pdf数据产生失败！")
		return false
	}
	// 发送微信
	resp := wecom.SendFileToUser(pdfBytes, ".pdf", username)
	if resp == nil {
		xlog.Log.WithField("用户:", username).WithError(err).Error("获取导出用户历史发送文件消息失败！")
		return false
	}
	return true
}

// SearchHistoryAndGenPDF
//  @Description: 生成导出的pdf Byte文件
//  @param context
//  @param err
//  @param username
//  @return []byte
//  @return error
//
func SearchHistoryAndGenPDF(context *dao.Context, err error, username string) ([]byte, error) {
	contentLastModTime := context.UpdateTime
	msgContext, err := UnMarshalJSonToMsgContext(context.Name, context.ContextMsg)
	if err != nil {
		xlog.Log.WithField("用户:", username).WithError(err).Error("导出历史记录系统错误！")
		return nil, err
	}
	buffer := bytes.Buffer{}
	// 写入文档标题
	buffer.Write([]byte(GetOutPutHeader(username, contentLastModTime)))
	// 写入对话
	for key, message := range msgContext.Context {
		buffer.Write([]byte(GetChatRoleHeader(message.Role, username, key)))
		buffer.Write([]byte(GetChatContentByRole(message.Content, message.Role)))
	}
	// 解析成html
	dataHtml, err := doc.GetHtmlFromMdBytes(GetGenDocTitle(username), buffer.Bytes())
	if err != nil {
		xlog.Log.WithField("用户:", username).WithError(err).Error("导出历史记录过程中解析生成的记录文件失败")
		return nil, err
	}
	// 转化pdf
	pdf := doc.ConvertHtmlToPDF([]byte(dataHtml))
	if pdf == nil {
		xlog.Log.WithField("用户:", username).Error("导出历史记录过程中pdf文件生成失败")
		return nil, errors.New("pdf工具类转化pdf失败")
	}
	return pdf, nil
}

// GetOutPutHeader
//
//	@Description: 生成输出文档大标题
//	@param content
//	@param timeStamp
//	@return string
func GetOutPutHeader(content string, timeStamp int64) string {
	t := time.UnixMilli(timeStamp).String()
	return "\n\n# " + content + "TIME-" + t
}

// GetChatRoleHeader
//
//	@Description: 获取角色对话标题
//	@param role
//	@param userName
//	@return string
func GetChatRoleHeader(role string, userName string, key int) string {
	who := "预定提示词"
	if key == 0 {
		return fmt.Sprintf("\n\n## No: %d %s said: ", key+1, who)
	}
	switch role {
	case model.SYSTEM:
		who = "预定提示词"
	case model.USER:
		who = userName
	case model.AI:
		who = "AI"
	default:
		who = "not set"
	}
	return fmt.Sprintf("\n\n## No: %d %s said: ", key+1, who)
}

// GetChatContentByRole
//  @Description: 获取格式化的聊天内容
//  @param content
//  @param role
//  @return string
//
func GetChatContentByRole(content string, role string) string {
	switch role {
	case model.AI:
		return content
	default:
		return "\n\n" + content
	}
}

// GetGenDocTitle
//  @Description: 获取文件标题
//  @param userName
//  @return string
//
func GetGenDocTitle(userName string) string {
	return userName + "对话内容"
}
