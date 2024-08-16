package authority

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sort"

	"server/global"
	modelAuthority "server/model/authority"
	authorityReq "server/model/authority/request"
)

type MenuService struct{}

func getTreeMap(menuListFormat []modelAuthority.MenuModel, menuList []modelAuthority.MenuModel) {
	for index, menuF := range menuListFormat {
		for _, menu := range menuList {
			if menuF.ID == menu.Pid {
				// menuF 只是个复制值
				//menuF.Children = append(menuF.Children, menu)
				menuListFormat[index].Children = append(menuListFormat[index].Children, menu)
			}
		}
		if len(menuListFormat[index].Children) > 0 {
			// 排序
			sort.Slice(menuListFormat[index].Children, func(i, j int) bool {
				return menuListFormat[index].Children[i].Sort < menuListFormat[index].Children[j].Sort
			})
			getTreeMap(menuListFormat[index].Children, menuList)
		}
	}
}

func (ms *MenuService) GetMenus(userId uint) ([]modelAuthority.MenuModel, error) {
	// 查找用户
	var userModel modelAuthority.UserModel
	err := global.Db.Where("id = ?", userId).First(&userModel).Error
	if err != nil {
		return nil, fmt.Errorf("GetMenus 用户查询 -> %v", err)
	}

	// 查找用户对应角色
	var roleModel modelAuthority.RoleModel
	err = global.Db.Where("id = ?", userModel.RoleModelID).First(&roleModel).Error
	if err != nil {
		return nil, fmt.Errorf("GetMenus 角色查询 -> %v", err)
	}

	var menuModels []modelAuthority.MenuModel
	err = global.Db.Preload("Roles").Find(&menuModels).Error
	if err != nil {
		return nil, fmt.Errorf("GetMenus 菜单查询 -> %v", err)
	}

	// 过滤角色拥有的路由
	menuList := make([]modelAuthority.MenuModel, 0)
	for _, menu := range menuModels {
		for _, menuRole := range menu.Roles {
			if roleModel.RoleName == menuRole.RoleName {
				menuList = append(menuList, menu)
				continue
			}
		}
	}

	// 找出第一级路由，（父路由id为0）
	menuListFormat := make([]modelAuthority.MenuModel, 0)
	for _, menu := range menuList {
		if menu.Pid == 0 {
			menuListFormat = append(menuListFormat, menu)
		}
	}

	// 排序
	sort.Slice(menuListFormat, func(i, j int) bool {
		return menuListFormat[i].Sort < menuListFormat[j].Sort
	})

	// 递归找出一级路由下面的子路由
	getTreeMap(menuListFormat, menuList)

	return menuListFormat, nil
}

func (ms *MenuService) AddMenu(menuRaw authorityReq.Menu) bool {
	var menuModel modelAuthority.MenuModel
	menuModel.Name = menuRaw.Name
	menuModel.Path = menuRaw.Path
	menuModel.Component = menuRaw.Component
	menuModel.Redirect = menuRaw.Redirect
	menuModel.Pid = menuRaw.Pid
	menuModel.Sort = menuRaw.Sort
	menuModel.Meta.Title = menuRaw.Meta.Title
	menuModel.Meta.SvgIcon = menuRaw.Meta.Icon
	menuModel.Meta.Hidden = menuRaw.Meta.Hidden
	menuModel.Meta.Affix = menuRaw.Meta.Affix
	menuModel.Meta.KeepAlive = menuRaw.Meta.KeepAlive
	menuModel.Meta.AlwaysShow = menuRaw.Meta.AlwaysShow

	if err := global.Db.Create(&menuModel).Error; err != nil {
		global.Log.Error("创建menu失败", zap.Error(err))
		return false
	}

	return true
}

func (ms *MenuService) EditMenu(menuRaw authorityReq.EditMenuReq) (err error) {
	var menuModel modelAuthority.MenuModel
	var metaData modelAuthority.Meta

	if errors.Is(global.Db.Where("id = ?", menuRaw.ID).First(&menuModel).Error, gorm.ErrRecordNotFound) {
		return errors.New("菜单不存在")
	}

	metaData.SvgIcon = menuRaw.Meta.Icon
	metaData.Title = menuRaw.Meta.Title
	metaData.Hidden = menuRaw.Meta.Hidden
	metaData.Affix = menuRaw.Meta.Affix
	metaData.KeepAlive = menuRaw.Meta.KeepAlive
	metaData.AlwaysShow = menuRaw.Meta.AlwaysShow

	err = global.Db.Model(&menuModel).Updates(map[string]interface{}{
		"pid":       menuRaw.Pid,
		"name":      menuRaw.Name,
		"path":      menuRaw.Path,
		"component": menuRaw.Component,
		"redirect":  menuRaw.Redirect,
		"sort":      menuRaw.Sort,
		"meta":      metaData,
	}).Error

	return
}

func (ms *MenuService) DeleteMenu(id uint) (err error) {
	var menuModel modelAuthority.MenuModel
	if errors.Is(global.Db.Where("id = ?", id).First(&menuModel).Error, gorm.ErrRecordNotFound) {
		return errors.New("菜单不存在")
	}
	err = global.Db.Unscoped().Select("Roles").Delete(&menuModel).Error

	return err
}

// GetElTreeMenus 获取所有menu
func (ms *MenuService) GetElTreeMenus(roleId uint) ([]modelAuthority.MenuModel, []uint, error) {
	var menuModels []modelAuthority.MenuModel
	err := global.Db.Find(&menuModels).Error
	if err != nil {
		global.Log.Error("GetElTreeMenus 查询menus", zap.Error(err))
		return nil, nil, err
	}

	menuListFormat := make([]modelAuthority.MenuModel, 0)
	for _, menu := range menuModels {
		if menu.Pid == 0 {
			menuListFormat = append(menuListFormat, menu)
		}
	}

	getTreeMap(menuListFormat, menuModels)

	var roleModel modelAuthority.RoleModel
	err = global.Db.Where("id = ?", roleId).Preload("Menus").First(&roleModel).Error
	if err != nil {
		global.Log.Error("GetElTreeMenus 查询role", zap.Error(err))
		return nil, nil, err
	}

	// 前端el-tree 选中数据
	// 去掉夫菜单，防止直接选中父级造成全选
	roleIds := make([]uint, 0)
	count := 0
	for _, menu := range roleModel.Menus {
		for _, menu1 := range roleModel.Menus {
			if menu.ID == menu1.Pid {
				count++
				break
			}
		}
		if count == 0 {
			roleIds = append(roleIds, menu.ID)
		} else {
			count--
		}
	}

	return menuListFormat, roleIds, nil
}
