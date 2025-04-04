package system

import (
	"eve-corp-manager/global"
	"eve-corp-manager/models/system"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) Add(user *system.User) (*system.User, error) {
	err := r.DB.Create(user).Error
	if err != nil {
		global.Logger.Errorf("Failed to add user, got error: %v", err)
		return nil, err
	}
	return user, err
}

func (r *UserRepository) Get(user *system.User) (*system.User, error) {
	err := r.DB.Where("user_id = ?", user.UserId).First(user).Error
	if err != nil {
		global.Logger.Errorf("Failed to get user, got error: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(user *system.User) (*system.User, error) {
	err := r.DB.Save(user).Error
	if err != nil {
		global.Logger.Errorf("Failed to update user, got error: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetCharacterList(userID uint) (*system.User, error) {
	var user system.User
	err := r.DB.Preload("Characters").First(&user, userID).Error
	if err != nil {
		global.Logger.Errorf("Failed to get user character list, got error: %v", err)
		return nil, err
	}
	return &user, nil
}
