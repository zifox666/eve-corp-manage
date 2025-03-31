package system

import (
	"eve-corp-manager/models/common"
)

type User struct {
	common.BaseModelNoId

	UserId          uint   `gorm:"index;type:uint"  json:"userId"`
	MainCharacterId int    `gorm:"index;type:int(11)" json:"mainCharacterId"` // EVE 主角色ID
	Qq              uint   `gorm:"type:int(11)" json:"qq"`                    // QQ号
	Name            string `gorm:"type:varchar(20)" json:"name"`              // 昵称
	Status          int    `gorm:"type:tinyint(1)" json:"status"`             // 用户状态
}
