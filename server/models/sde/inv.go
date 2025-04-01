package sde

import "eve-corp-manager/global"

// 翻译相关的常量
const (
	CategoryTcID = 6 // 物品分类的翻译ID
	GroupTcID    = 7 // 物品组的翻译ID
	NameTcID     = 8 // 物品名称的翻译ID
)

// InvFlag EVE游戏中物品的标志定义
type InvFlag struct {
	FlagID   int    `gorm:"primaryKey;column:flagID"`
	FlagName string `gorm:"column:flagName"`
	FlagText string `gorm:"column:flagText"`
	OrderID  int    `gorm:"column:orderID"`
}

// TableName 指定表名
func (InvFlag) TableName() string {
	return "invFlags"
}

// InvGroup 物品组定义
type InvGroup struct {
	GroupID              int    `gorm:"primaryKey;column:groupID"`
	CategoryID           int    `gorm:"column:categoryID"`
	GroupName            string `gorm:"column:groupName"`
	IconID               int    `gorm:"column:iconID"`
	UseBasePrice         bool   `gorm:"column:useBasePrice"`
	Anchored             bool   `gorm:"column:anchored"`
	Anchorable           bool   `gorm:"column:anchorable"`
	FittableNonSingleton bool   `gorm:"column:fittableNonSingleton"`
	Published            bool   `gorm:"column:published"`
}

// TableName 指定表名
func (InvGroup) TableName() string {
	return "invGroups"
}

// InvType 物品类型定义
type InvType struct {
	TypeID        int     `gorm:"primaryKey;column:typeID"`
	GroupID       int     `gorm:"column:groupID"`
	TypeName      string  `gorm:"column:typeName"`
	Description   string  `gorm:"column:description"`
	Mass          float64 `gorm:"column:mass"`
	Volume        float64 `gorm:"column:volume"`
	Capacity      float64 `gorm:"column:capacity"`
	PortionSize   int     `gorm:"column:portionSize"`
	RaceID        int     `gorm:"column:raceID"`
	BasePrice     float64 `gorm:"column:basePrice;type:decimal(19,4)"`
	Published     bool    `gorm:"column:published"`
	MarketGroupID int     `gorm:"column:marketGroupID"`
	IconID        int     `gorm:"column:iconID"`
	SoundID       int     `gorm:"column:soundID"`
	GraphicID     int     `gorm:"column:graphicID"`
}

// TableName 指定表名
func (InvType) TableName() string {
	return "invTypes"
}

// TrnTranslation 多语言翻译表
type TrnTranslation struct {
	TcID       int    `gorm:"primaryKey;column:tcID"`
	KeyID      int    `gorm:"primaryKey;column:keyID"`
	LanguageID string `gorm:"primaryKey;column:languageID"`
	Text       string `gorm:"column:text"`
}

// TableName 指定表名
func (TrnTranslation) TableName() string {
	return "trnTranslations"
}

// TypeInfo 包含物品的ID层级和对应名称
type TypeInfo struct {
	TypeID       int    `json:"typeId"`
	TypeName     string `json:"typeName"`
	GroupID      int    `json:"groupId"`
	GroupName    string `json:"groupName"`
	CategoryID   int    `json:"categoryId"`
	CategoryName string `json:"categoryName,omitempty"`
}

// GetTypeInfoByID 根据物品ID和语言代码获取物品及其分类信息
func GetTypeInfoByID(typeID int, lang string) (TypeInfo, error) {
	var info TypeInfo

	// 1. 获取基本信息（type -> group -> category 的ID链和英文名称）
	err := global.SdeDb.Raw(`
		SELECT t.typeID, t.typeName, t.groupID, g.groupName, g.categoryID
		FROM invTypes t
		JOIN invGroups g ON t.groupID = g.groupID
		WHERE t.typeID = ?
	`, typeID).Scan(&info).Error

	if err != nil {
		return info, err
	}

	// 2. 如果指定了语言，获取对应的翻译名称
	if lang != "" {
		// 查询物品名称翻译
		var typeTrans TrnTranslation
		global.SdeDb.Where("tcID = ? AND keyID = ? AND languageID = ?",
			NameTcID, typeID, lang).First(&typeTrans)
		if typeTrans.Text != "" {
			info.TypeName = typeTrans.Text
		}

		// 查询分组名称翻译
		var groupTrans TrnTranslation
		global.SdeDb.Where("tcID = ? AND keyID = ? AND languageID = ?",
			GroupTcID, info.GroupID, lang).First(&groupTrans)
		if groupTrans.Text != "" {
			info.GroupName = groupTrans.Text
		}

		// 查询分类名称翻译
		var categoryTrans TrnTranslation
		global.SdeDb.Where("tcID = ? AND keyID = ? AND languageID = ?",
			CategoryTcID, info.CategoryID, lang).First(&categoryTrans)
		if categoryTrans.Text != "" {
			info.CategoryName = categoryTrans.Text
		}
	}

	return info, nil
}
