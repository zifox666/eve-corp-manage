package service

import (
	"eve-corp-manager/models/common"
	"time"
)

type Fleet struct {
	common.BaseModel

	FleetName          string                 `gorm:"type:varchar(255)" json:"fleetName"`      // 舰队名称
	FleetType          int                    `gorm:"type:tinyint(1)" json:"fleetType"`        // 舰队类型
	FleetExtraInfo     string                 `gorm:"type:varchar(255)" json:"fleetInfo"`      // 舰队自定义信息
	FleetCommanderID   int                    `gorm:"type:int(11)" json:"fleetCommanderId"`    // 舰队指挥ID
	FleetCommanderName string                 `gorm:"type:varchar(255)" json:"fleetCommander"` // 舰队指挥名称
	FleetLocation      map[string]interface{} `gorm:"type:json" json:"fleetLocation"`          // 舰队位置集合
	Sig                int                    `gorm:"type:int(8);default:0" json:"sig"`        // 小组类别
	Srp                bool                   `gorm:"type:tinyint(1)" json:"srp"`              // 是否有SRP
	CorpPap            float32                `gorm:"type:float" json:"corpPap"`               // 公司PAP
	AutoSrp            bool                   `gorm:"type:tinyint(1)" json:"autoSrp"`          // 自动SRP
	StartTime          time.Time
	EndTime            time.Time

	CharacterFleetAssociations []CharacterFleetAssociation `gorm:"foreignKey:FleetID;references:FleetID"`
}

type CharacterFleetAssociation struct {
	common.BaseModel

	FleetID     uint `gorm:"index;type:uint" json:"fleetId"`     // 舰队ID
	CharacterID uint `gorm:"index;type:uint" json:"characterId"` // 角色ID

	Fleet Fleet `gorm:"foreignKey:FleetID;references:FleetID"`
}
