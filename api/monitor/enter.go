package monitor

import "server/service"

type ApiGroup struct {
	OperationLogApi
}

var (
	operationLogService = service.GroupApp.Monitor.OperationLogService
)
