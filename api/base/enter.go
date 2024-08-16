package base

import (
	"server/service"
)

type ApiGroup struct {
	LogRegApi
	CasbinApi
}

var (
	jwtService    = service.GroupApp.Base.JwtService
	logRegService = service.GroupApp.Base.LogRegService
	casbinService = service.GroupApp.Base.CasbinService
)
