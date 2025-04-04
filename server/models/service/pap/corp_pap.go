package pap

import (
	"eve-corp-manager/models/common"
	"time"
)

// CorpPap 用户PAP记录表
type CorpPap struct {
	common.BaseModel
	UserID     uint      `gorm:"index;type:uint" json:"userId"`   // 用户ID
	Amount     int       `gorm:"type:int" json:"amount"`          // PAP数量
	Balance    int       `gorm:"type:int" json:"balance"`         // PAP余额
	Source     string    `gorm:"type:varchar(255)" json:"source"` // PAP来源
	SourceID   uint      `gorm:"type:uint" json:"sourceId"`       // 来源ID（如舰队ID）
	Type       int       `gorm:"type:tinyint(1)" json:"type"`     // 类型：1-获取 2-消费
	CreateTime time.Time `gorm:"type:datetime" json:"createTime"` // 创建时间
	Remark     string    `gorm:"type:varchar(255)" json:"remark"` // 备注
}

// CorpPapLog PAP操作日志表
type CorpPapLog struct {
	common.BaseModel
	UserID     uint      `gorm:"index;type:uint" json:"userId"`     // 用户ID
	Operation  string    `gorm:"type:varchar(50)" json:"operation"` // 操作类型
	Amount     int       `gorm:"type:int" json:"amount"`            // 操作数量
	BeforeVal  int       `gorm:"type:int" json:"beforeVal"`         // 操作前值
	AfterVal   int       `gorm:"type:int" json:"afterVal"`          // 操作后值
	Operator   uint      `gorm:"type:uint" json:"operator"`         // 操作人ID
	CreateTime time.Time `gorm:"type:datetime" json:"createTime"`   // 创建时间
	Remark     string    `gorm:"type:varchar(255)" json:"remark"`   // 备注
}

// CorpPapShopItem 商城兑换项目表（预留）
type CorpPapShopItem struct {
	common.BaseModel
	ItemName   string    `gorm:"type:varchar(100)" json:"itemName"` // 商品名称
	ItemDesc   string    `gorm:"type:varchar(255)" json:"itemDesc"` // 商品描述
	PapCost    int       `gorm:"type:int" json:"papCost"`           // PAP消耗
	Status     int       `gorm:"type:tinyint(1)" json:"status"`     // 状态：1-可用 0-不可用
	CreateTime time.Time `gorm:"type:datetime" json:"createTime"`   // 创建时间
	UpdateTime time.Time `gorm:"type:datetime" json:"updateTime"`   // 更新时间
	Remark     string    `gorm:"type:varchar(255)" json:"remark"`   // 备注
}
