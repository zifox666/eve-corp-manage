package character

import (
	"eve-corp-manager/global"
	"eve-corp-manager/models/service/character"
	"eve-corp-manager/utils"
	"gorm.io/gorm"
)

type UserCharacterRepository struct {
	DB *gorm.DB
}

func (r *UserCharacterRepository) Update(userID, characterID uint, newUserCharacter character.UserCharacter) (*character.UserCharacter, error) {
	err := r.DB.Model(&character.UserCharacter{}).
		Where("user_id = ? AND character_id = ?", userID, characterID).
		Save(newUserCharacter).
		Error
	if err != nil {
		global.Logger.Errorf("Failed to update user character, userID: %v, characterID: %v, error: %v", userID, characterID, err)
		return nil, err
	}
	return &newUserCharacter, nil
}

func (r *UserCharacterRepository) GetAllInAllowedCorp() ([]character.UserCharacter, error) {
	corpList, err := global.Settings.Get("allowed_corp_list")
	if err != nil {
		global.Logger.Errorf("allowed_corp_list变量未设置，提取所有\n %v", err)
	}
	corpIdList, err := utils.StringToIntList(corpList)
	result := r.DB.Where("corporation_id IN ?", corpIdList).Find(&character.UserCharacter{})
	if result.Error != nil {
		global.Logger.Errorf("获取公司列表失败: %v", result.Error)
		return nil, result.Error
	}
	return result, nil
}
