package service

import (
	"eve-corp-manager/models/common"
	"time"
)

// Fleet 结构体
type Fleet struct {
	common.BaseModel
	FleetName          string
	FleetType          int
	FleetExtraInfo     string
	FleetCommanderID   uint
	FleetCommanderName string
	FleetLocation      map[string]interface{} `gorm:"serializer:json"`
	Sig                int
	Srp                bool
	CorpPap            int
	AutoSrp            bool
	StartTime          time.Time
	EndTime            time.Time
	// 修改关联关系定义
	CharacterFleetAssociations []CharacterFleetAssociation `gorm:"foreignKey:FleetID"`
}

// CharacterFleetAssociation 结构体
type CharacterFleetAssociation struct {
	common.BaseModel
	FleetID     uint `gorm:"index"` // 外键字段
	CharacterID uint `gorm:"index"` // 角色ID
	// 修改关联关系定义
	Fleet Fleet `gorm:"foreignKey:FleetID"` // 使用FleetID作为外键
}
