package services

import (
	"errors"

	"corp-webot/xconst"
	"golang.org/x/net/context"
)

var commandSupported map[string]CommonMsgCmd

func init() {
	// 初始化命令
	commandSupported = make(map[string]CommonMsgCmd, 10)
	// 注册命令
	RegisterCommand(xconst.CPT_CMD, &GPTChatCommand{})
}

type CommandData struct {
	CorpID   string
	FromUser string
	Msg      string
	Cmd      string
}

func NewCommandData(cropId, fromUser, msg, cmd string) *CommandData {
	return &CommandData{
		CorpID:   cropId,
		FromUser: fromUser,
		Msg:      msg,
		Cmd:      cmd,
	}
}

type CommonMsgCmd interface {
	ExecCommand(data *CommandData, ctx context.Context)
}

func RegisterCommand(cmdName string, cmd CommonMsgCmd) {
	commandSupported[cmdName] = cmd
}

func GetCommandFunc(cmdName string) (CommonMsgCmd, error) {
	cmd, ok := commandSupported[cmdName]
	if !ok {
		return nil, errors.New("未找到指定文本命令")
	}
	return cmd, nil
}
