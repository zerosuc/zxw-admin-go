package initialize

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math/rand"
	"server/global"
	"server/model/authority"
	"server/model/fileM"
	"server/model/sysTool"
	"strconv"
)

func InitData(db *gorm.DB) {
	initApiData(db)
	//InitMenuData(db)
	initRoleData(db)
	initUserData(db)
	initCasbinData(db)
	initFileMData(db)
	initCronData(db)
}
func initApiData(db *gorm.DB) {
	// 定义要插入的数据
	userEntries := []authority.ApiModel{
		{Path: "/logReg/captcha", Description: "获取验证码（必选）", ApiGroup: "logReg", Method: "POST"},
		{Path: "/logReg/login", Description: "登录（必选）", ApiGroup: "logReg", Method: "POST"},
		{Path: "/logReg/logout", Description: "登出（必选）", ApiGroup: "logReg", Method: "POST"},
		{Path: "/casbin/editCasbin", Description: "编辑casbin规则", ApiGroup: "casbin", Method: "POST"},
		{Path: "/user/getUserInfo", Description: "获取用户信息（必选）", ApiGroup: "user", Method: "GET"},
		{Path: "/user/getUsers", Description: "获取所有用户", ApiGroup: "user", Method: "POST"},
		{Path: "/user/deleteUser", Description: "删除用户", ApiGroup: "user", Method: "POST"},
		{Path: "/user/addUser", Description: "添加用户", ApiGroup: "user", Method: "POST"},
		{Path: "/user/editUser", Description: "编辑用户", ApiGroup: "user", Method: "POST"},
		{Path: "/user/modifyPass", Description: "修改用户密码", ApiGroup: "user", Method: "POST"},
		{Path: "/user/switchActive", Description: "切换用户状态", ApiGroup: "user", Method: "POST"},
		{Path: "/role/getRoles", Description: "获取所有角色", ApiGroup: "role", Method: "POST"},
		{Path: "/role/addRole", Description: "添加角色", ApiGroup: "role", Method: "POST"},
		{Path: "/role/deleteRole", Description: "删除角色", ApiGroup: "role", Method: "POST"},
		{Path: "/role/editRole", Description: "编辑角色", ApiGroup: "role", Method: "POST"},
		{Path: "/role/editRoleMenu", Description: "编辑角色菜单", ApiGroup: "role", Method: "POST"},
		{Path: "/menu/getMenus", Description: "获取所有菜单", ApiGroup: "menu", Method: "GET"},
		{Path: "/menu/addMenu", Description: "添加菜单", ApiGroup: "menu", Method: "POST"},
		{Path: "/menu/editMenu", Description: "编辑菜单", ApiGroup: "menu", Method: "POST"},
		{Path: "/menu/deleteMenu", Description: "删除菜单", ApiGroup: "menu", Method: "POST"},
		{Path: "/menu/getElTreeMenus", Description: "获取所有菜单（el-tree结构）", ApiGroup: "menu", Method: "POST"},
		{Path: "/api/addApi", Description: "添加api", ApiGroup: "api", Method: "POST"},
		{Path: "/api/getApis", Description: "获取所有api", ApiGroup: "api", Method: "POST"},
		{Path: "/api/deleteApi", Description: "删除api", ApiGroup: "api", Method: "POST"},
		{Path: "/api/editApi", Description: "编辑api", ApiGroup: "api", Method: "POST"},
		{Path: "/api/getElTreeApis", Description: "获取所有api（el-tree结构）", ApiGroup: "api", Method: "POST"},
		{Path: "/api/deleteApiById", Description: "批量删除API", ApiGroup: "api", Method: "POST"},
		{Path: "/opl/getOplList", Description: "分页获取操作记录", ApiGroup: "opl", Method: "POST"},
		{Path: "/opl/deleteOpl", Description: "删除操作记录", ApiGroup: "opl", Method: "POST"},
		{Path: "/opl/deleteOplByIds", Description: "批量删除操作记录", ApiGroup: "opl", Method: "POST"},
		{Path: "/file/upload", Description: "文件上传", ApiGroup: "file", Method: "POST"},
		{Path: "/file/getFileList", Description: "分页获取文件信息", ApiGroup: "file", Method: "POST"},
		{Path: "/file/download", Description: "下载文件", ApiGroup: "file", Method: "GET"},
		{Path: "/file/delete", Description: "删除文件", ApiGroup: "file", Method: "GET"},
		{Path: "/cron/getCronList", Description: "分页获取cron", ApiGroup: "cron", Method: "POST"},
		{Path: "/cron/addCron", Description: "添加cron", ApiGroup: "cron", Method: "POST"},
		{Path: "/cron/deleteCron", Description: "删除cron", ApiGroup: "cron", Method: "POST"},
		{Path: "/cron/editCron", Description: "编辑cron", ApiGroup: "cron", Method: "POST"},
		{Path: "/cron/switchOpen", Description: "cron开关", ApiGroup: "cron", Method: "POST"},
	}

	// 使用 GORM 插入数据
	if err := db.Create(&userEntries).Error; err != nil {
		global.Log.Error("router register ApiModel failed")
	} else {
		global.Log.Info("router register ApiModel success")
	}
}

func initMenuData(db *gorm.DB) {
	menuEntries := []authority.MenuModel{
		{
			Pid:       0,
			Name:      "Authority",
			Path:      "/authority",
			Redirect:  "/authority/user",
			Component: "Layout",
			Meta:      authority.Meta{Title: "权限管理", SvgIcon: "lock"},
			Sort:      1,
		},
		{
			Pid:       1,
			Name:      "User",
			Path:      "user",
			Component: "authority/user/index.vue",
			Meta:      authority.Meta{Title: "用户管理"},
			Sort:      1,
		},
		{
			Pid:       1,
			Name:      "Role",
			Path:      "role",
			Component: "authority/role/index.vue",
			Meta:      authority.Meta{Title: "角色管理"},
			Sort:      2,
		},
		{
			Pid:       1,
			Name:      "Menu",
			Path:      "menu",
			Component: "authority/menu/index.vue",
			Meta:      authority.Meta{Title: "菜单管理"},
			Sort:      3,
		},
		{
			Pid:       1,
			Name:      "Api",
			Path:      "api",
			Component: "authority/api/index.vue",
			Meta:      authority.Meta{Title: "接口管理"},
			Sort:      4,
		},
		{
			Pid:       0,
			Name:      "Cenu",
			Path:      "/cenu",
			Redirect:  "/cenu/cenu1",
			Component: "Layout",
			Meta:      authority.Meta{Title: "多级菜单", SvgIcon: "menu", AlwaysShow: true},
			Sort:      5,
		},
		{
			Pid:       6,
			Name:      "Cenu1",
			Path:      "cenu1",
			Redirect:  "/cenu/cenu1/cenu1-1",
			Component: "cenu/cenu1/index.vue",
			Meta:      authority.Meta{Title: "cenu1"},
			Sort:      1,
		},
		{
			Pid:       7,
			Name:      "Cenu1-1",
			Path:      "cenu1-1",
			Component: "cenu/cenu1/cenu1-1/index.vue",
			Meta:      authority.Meta{Title: "cenu1-1"},
			Sort:      1,
		},
		{
			Pid:       7,
			Name:      "Cenu1-2",
			Path:      "cenu1-2",
			Component: "cenu/cenu1/cenu1-2/index.vue",
			Meta:      authority.Meta{Title: "cenu1-2"},
			Sort:      2,
		},
		{
			Pid:       14,
			Name:      "File",
			Path:      "/fileM/file",
			Component: "fileM/file/index.vue",
			Meta:      authority.Meta{Title: "文件上传"},
			Sort:      1,
		},
		{
			Pid:       14,
			Name:      "OperationLog",
			Path:      "operationLog",
			Component: "monitor/operationLog/index.vue",
			Meta:      authority.Meta{Title: "操作日志"},
			Sort:      1,
		},
		{
			Pid:       0,
			Name:      "SysTool",
			Path:      "/systool",
			Redirect:  "/systool/cron",
			Component: "Layout",
			Meta:      authority.Meta{Title: "系统工具", SvgIcon: "config", AlwaysShow: true},
			Sort:      4,
		},
		{
			Pid:       14,
			Name:      "Cron",
			Path:      "cron",
			Component: "sysTool/cron/index.vue",
			Meta:      authority.Meta{Title: "定时任务"},
			Sort:      1,
		},
	}
	if err := db.Create(&menuEntries).Error; err != nil {
		global.Log.Error("router register menuEntries failed")
	} else {
		global.Log.Info("router register menuEntries success")
	}
}

func initUserData(db *gorm.DB) {
	userEntries := []authority.UserModel{
		{
			Username:    "admin",
			Password:    "e10adc3949ba59abbe56e057f20f883e",
			Phone:       "16666666666",
			Email:       "wanminny@163.com",
			Active:      true,
			RoleModelID: 1,
		},
		{
			Username:    "zxw",
			Password:    "e10adc3949ba59abbe56e057f20f883e",
			Phone:       "13333333333",
			Email:       "sfasfd@fad.com",
			Active:      true,
			RoleModelID: 2,
		},
	}
	if err := db.Create(&userEntries).Error; err != nil {
		global.Log.Error("router register userEntries failed")
	} else {
		global.Log.Info("router register userEntries success")
	}
}
func initFileMData(db *gorm.DB) {
	file1 := []fileM.FileModel{
		{
			FileName: "k8s_4687661d-bea5-4883-ad63-1112263355cb.csv",
			FullPath: "./resource/upload/k8s_4687661d-bea5-4883-ad63-1112263355cb.csv",
			Mime:     "text/csv",
		},
		{
			FileName: "hosts_91fd2972-b38b-45af-a0fd-17d9491cb3f7.csv",
			FullPath: "./resource/upload/hosts_91fd2972-b38b-45af-a0fd-17d9491cb3f7.csv",
			Mime:     "text/csv",
		},
	}
	if err := db.Create(&file1).Error; err != nil {
		global.Log.Error("router register file1 failed")
	} else {
		global.Log.Info("router register file1 success")
	}

}

func initCasbinData(db *gorm.DB) {
	casbinRules := []gormadapter.CasbinRule{
		{Ptype: "p", V0: "1", V1: "/api/addApi", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/api/deleteApi", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/api/deleteApiById", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/api/editApi", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/api/getApis", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/api/getElTreeApis", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/casbin/editCasbin", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/cron/addCron", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/cron/deleteCron", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/cron/editCron", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/cron/getCronList", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/cron/switchOpen", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/file/delete", V2: "GET"},
		{Ptype: "p", V0: "1", V1: "/file/download", V2: "GET"},
		{Ptype: "p", V0: "1", V1: "/file/getFileList", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/file/upload", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/logReg/captcha", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/logReg/login", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/logReg/logout", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/menu/addMenu", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/menu/deleteMenu", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/menu/editMenu", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/menu/getElTreeMenus", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/menu/getMenus", V2: "GET"},
		{Ptype: "p", V0: "1", V1: "/opl/deleteOpl", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/opl/deleteOplByIds", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/opl/getOplList", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/role/addRole", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/role/deleteRole", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/role/editRole", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/role/editRoleMenu", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/role/getRoles", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/user/addUser", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/user/deleteUser", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/user/editUser", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/user/getUserInfo", V2: "GET"},
		{Ptype: "p", V0: "1", V1: "/user/getUsers", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/user/modifyPass", V2: "POST"},
		{Ptype: "p", V0: "1", V1: "/user/switchActive", V2: "POST"},
		{Ptype: "p", V0: "2", V1: "/logReg/captcha", V2: "POST"},
		{Ptype: "p", V0: "2", V1: "/logReg/login", V2: "POST"},
		{Ptype: "p", V0: "2", V1: "/logReg/logout", V2: "POST"},
		{Ptype: "p", V0: "2", V1: "/menu/getMenus", V2: "GET"},
		{Ptype: "p", V0: "2", V1: "/user/editUser", V2: "POST"},
		{Ptype: "p", V0: "2", V1: "/user/getUserInfo", V2: "GET"},
		{Ptype: "p", V0: "2", V1: "/user/modifyPass", V2: "POST"},
	}
	if err := db.Create(&casbinRules).Error; err != nil {
		global.Log.Error("router register casbinRules failed")
	} else {
		global.Log.Info("router register casbinRules success")
	}
}

// InitRoleData 关联表 role & menu
func initRoleData(db *gorm.DB) {
	roleEntries := []authority.RoleModel{
		{
			RoleName: "super",
			Menus: []*authority.MenuModel{
				&authority.MenuModel{
					Pid:       0,
					Name:      "Authority",
					Path:      "/authority",
					Redirect:  "/authority/user",
					Component: "Layout",
					Meta:      authority.Meta{Title: "权限管理", SvgIcon: "lock"},
					Sort:      1,
				},
				&authority.MenuModel{
					Pid:       1,
					Name:      "User",
					Path:      "user",
					Component: "authority/user/index.vue",
					Meta:      authority.Meta{Title: "用户管理"},
					Sort:      1,
				},
				&authority.MenuModel{
					Pid:       1,
					Name:      "Role",
					Path:      "role",
					Component: "authority/role/index.vue",
					Meta:      authority.Meta{Title: "角色管理"},
					Sort:      2,
				},
				&authority.MenuModel{
					Pid:       1,
					Name:      "Menu",
					Path:      "menu",
					Component: "authority/menu/index.vue",
					Meta:      authority.Meta{Title: "菜单管理"},
					Sort:      3,
				},
				&authority.MenuModel{
					Pid:       1,
					Name:      "Api",
					Path:      "api",
					Component: "authority/api/index.vue",
					Meta:      authority.Meta{Title: "接口管理"},
					Sort:      4,
				},
				&authority.MenuModel{
					Pid:       0,
					Name:      "Cenu",
					Path:      "/cenu",
					Redirect:  "/cenu/cenu1",
					Component: "Layout",
					Meta:      authority.Meta{Title: "k8s集群管理", SvgIcon: "menu", AlwaysShow: true},
					Sort:      5,
				},
				&authority.MenuModel{
					Pid:       6,
					Name:      "Cenu1",
					Path:      "cenu1",
					Redirect:  "/cenu/cenu1/cenu1-1",
					Component: "cenu/cenu1/index.vue",
					Meta:      authority.Meta{Title: "工作负载"},
					Sort:      1,
				},
				&authority.MenuModel{
					Pid:       7,
					Name:      "Cenu1-1",
					Path:      "cenu1-1",
					Component: "cenu/cenu1/cenu1-1/index.vue",
					Meta:      authority.Meta{Title: "Deployment"},
					Sort:      1,
				},
				&authority.MenuModel{
					Pid:       7,
					Name:      "Cenu1-2",
					Path:      "cenu1-2",
					Component: "cenu/cenu1/cenu1-2/index.vue",
					Meta:      authority.Meta{Title: "DaemonSet"},
					Sort:      2,
				},
				&authority.MenuModel{
					Pid:       12,
					Name:      "File",
					Path:      "/fileM/file",
					Component: "fileM/file/index.vue",
					Meta:      authority.Meta{Title: "文件上传"},
					Sort:      1,
				},
				&authority.MenuModel{
					Pid:       12,
					Name:      "OperationLog",
					Path:      "operationLog",
					Component: "monitor/operationLog/index.vue",
					Meta:      authority.Meta{Title: "操作历史"},
					Sort:      1,
				},
				&authority.MenuModel{
					Pid:       0,
					Name:      "SysTool",
					Path:      "/systool",
					Redirect:  "/systool/cron",
					Component: "Layout",
					Meta:      authority.Meta{Title: "系统工具", SvgIcon: "config", AlwaysShow: true},
					Sort:      4,
				},
				&authority.MenuModel{
					Pid:       12,
					Name:      "Cron",
					Path:      "cron",
					Component: "sysTool/cron/index.vue",
					Meta:      authority.Meta{Title: "定时任务"},
					Sort:      1,
				},
			},
		},
		{
			RoleName: "普通用户",
		},
	}
	if err := db.Save(&roleEntries).Error; err != nil {
		global.Log.Error("router register roleEntries failed", zap.String("create", err.Error()))
	} else {
		global.Log.Info("router register roleEntries success")
	}
}

func initCronData(db *gorm.DB) {
	CronEnties := sysTool.CronModel{
		//BaseModel: global.BaseModel{
		//	//ID:        1,
		//	CreatedAt: time.Date(2024, 8, 16, 0, 48, 7, 0, time.UTC),
		//	UpdatedAt: time.Date(2024, 8, 16, 0, 50, 0, 0, time.UTC),
		//},
		Name:       "定时任务1" + strconv.Itoa(rand.Intn(10000)),
		Method:     "clearTable",
		Expression: "0 0/2 * * * ? ",
		Strategy:   "once",
		Open:       false,
		ExtraParams: sysTool.ExtraParams{
			TableInfo: []sysTool.ClearTable{
				{
					Interval:     "2",
					TableName:    "monitor_operation_log",
					CompareField: "created_at",
				},
			},
		},
		EntryId: 0,
		Comment: "测试执行",
	}

	if err := db.Create(&CronEnties).Error; err != nil {
		global.Log.Error("router register initCronData failed")
	} else {
		global.Log.Info("router register initCronData success")
	}
}
