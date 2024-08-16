package authority

import "server/global"

type RoleModel struct {
	global.BaseModel
	RoleName string `json:"roleName" gorm:"unique" binding:"required"`
	//Users    []*UserModel `json:"users"`
	Menus []*MenuModel `json:"menus" gorm:"many2many:role_menus;"`
}

func (RoleModel) TableName() string {
	return "authority_role"
}
