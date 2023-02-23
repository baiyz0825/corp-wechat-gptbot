package events

import (
	"context"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/contract"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
	"github.com/pkg/errors"
)

var eventServices map[string]EventService

func init() {
	// 初始化命令
	eventServices = make(map[string]EventService, 10)
	// 注册命令
	RegisterCommand(models.CALLBACK_EVENT_ENTER_AGENT, &NorMalEventService{})
}

type EventService interface {
	DealEvent(event contract.EventInterface, ctx context.Context)
}

func RegisterCommand(cmdName string, cmd EventService) {
	eventServices[cmdName] = cmd
}

func GetCommandFunc(cmdName string) (EventService, error) {
	cmd, ok := eventServices[cmdName]
	if !ok {
		return nil, errors.New("未找到指定事件处理器")
	}
	return cmd, nil
}
