package character

import "eve-corp-manager/models/common"

type UserCharacter struct {
	common.BaseModelNoId

	CharacterID   uint    `gorm:"primaryKey;type:uint" json:"character_id"`
	CharacterName string  `gorm:"type:varchar(50)" json:"character_name"`
	RefreshToken  string  `gorm:"type:varchar(255)" json:"refresh_token"`
	SkillPoint    float64 `gorm:"type:decimal(15,2)" json:"skill_point"`
	Isk           float64 `gorm:"type:decimal(15,2)" json:"isk"`
	Status        int     `gorm:"type:tinyint(1)" json:"status"`
	CorpID        uint    `gorm:"type:uint" json:"corp_id"`
	AllianceID    uint    `gorm:"type:uint" json:"alliance_id"`
}
