package character

import (
	"eve-corp-manager/models/common"
	"time"
)

// KillmailList 击毁邮件列表
type KillmailList struct {
	common.BaseModel
	KillMailID      int       `gorm:"column:kill_mail_id;type:int;index" json:"killMailId"`              // 击毁邮件ID
	KillMailHash    string    `gorm:"column:kill_mail_hash;type:varchar(64)" json:"killMailHash"`        // 击毁邮件哈希
	KillMailTime    time.Time `gorm:"column:kill_mail_time;type:datetime" json:"killMailTime"`           // 击毁时间
	SolarSystemID   int       `gorm:"column:solar_system_id;type:int" json:"solarSystemId"`              // 星系ID
	SolarSystemName string    `gorm:"column:solar_system_name;type:varchar(100)" json:"solarSystemName"` // 星系名称
	ShipTypeID      int       `gorm:"column:ship_type_id;type:int" json:"shipTypeId"`                    // 舰船类型ID
	ShipTypeName    string    `gorm:"column:ship_type_name;type:varchar(100)" json:"shipTypeName"`       // 舰船类型名称
	CharacterID     int       `gorm:"column:character_id;type:int;index" json:"characterId"`             // 角色ID
	CharacterName   string    `gorm:"column:character_name;type:varchar(100)" json:"characterName"`      // 角色名称
	CorporationID   int       `gorm:"column:corporation_id;type:int" json:"corporationId"`               // 公司ID
	CorporationName string    `gorm:"column:corporation_name;type:varchar(100)" json:"corporationName"`  // 公司名称
	AllianceID      int       `gorm:"column:alliance_id;type:int" json:"allianceId"`                     // 联盟ID
	AllianceName    string    `gorm:"column:alliance_name;type:varchar(100)" json:"allianceName"`        // 联盟名称
	UserID          uint      `gorm:"column:user_id;type:uint" json:"userId"`                            // 用户ID
	JaniceAmount    float64   `gorm:"column:janice_amount;type:decimal(20,2)" json:"janiceAmount"`       // Janice估价金额
	CreateBy        string    `gorm:"column:create_by;type:varchar(64)" json:"createBy"`                 // 创建者
	CreateID        uint      `gorm:"column:create_id;type:uint" json:"createId"`                        // 创建者ID
	// 关联关系
	Items []KillmailItem `gorm:"foreignKey:KillMailID;references:KillMailID" json:"items"` // 击毁邮件物品列表
}

// TableName 设置表名
func (KillmailList) TableName() string {
	return "killmail_list"
}

// KillmailItem 击毁邮件物品
type KillmailItem struct {
	common.BaseModel
	KillMailID int    `gorm:"column:kill_mail_id;type:int;index" json:"killMailId"` // 击毁邮件ID
	ItemID     int    `gorm:"column:item_id;type:int" json:"itemId"`                // 物品ID
	ItemName   string `gorm:"column:item_name;type:varchar(100)" json:"itemName"`   // 物品名称
	ItemNum    int    `gorm:"column:item_num;type:int" json:"itemNum"`              // 物品数量
	DropType   bool   `gorm:"column:drop_type;type:tinyint(1)" json:"dropType"`     // 掉落类型：true-掉落 false-摧毁
	SlotType   int    `gorm:"column:slot_type;type:int" json:"slotType"`            // 槽位类型
	CreateBy   string `gorm:"column:create_by;type:varchar(64)" json:"createBy"`    // 创建者
	CreateID   uint   `gorm:"column:create_id;type:uint" json:"createId"`           // 创建者ID
}

// TableName 设置表名
func (KillmailItem) TableName() string {
	return "killmail_item"
}
