package authority

import (
	"server/service"
)

type ApiGroup struct {
	UserApi
	MenuApi
	RoleApi
	ApiApi
}

var (
	userService   = service.GroupApp.Authority.UserService
	menuService   = service.GroupApp.Authority.MenuService
	roleService   = service.GroupApp.Authority.RoleService
	apiService    = service.GroupApp.Authority.ApiService
	casbinService = service.GroupApp.Base.CasbinService
)
