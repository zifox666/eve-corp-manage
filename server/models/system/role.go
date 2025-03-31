package system

import "eve-corp-manager/models/common"

type Role struct {
	common.BaseModel

	UserId uint `gorm:"primaryKey;index;type:uint" json:"userId"` // 用户ID
	RoleId uint `gorm:"primaryKey;index;type:uint" json:"roleId"` // 角色ID
}

type RoleMenu struct {
	common.BaseModel

	RoleId   uint `gorm:"primaryKey;index;type:uint" json:"roleId"` // 角色ID
	MenuId   uint `gorm:"primaryKey;index;type:uint" json:"menuId"` // 菜单ID
	IsButton bool `gorm:"type:tinyint(1)" json:"isButton"`          // 是否按钮
}
