package system

import (
	"encoding/json"
	"eve-corp-manager/core/cache"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SystemSetting 系统设置模型
type SystemSetting struct {
	gorm.Model
	ID          uint   `gorm:"primarykey"`
	ConfigName  string `gorm:"size:100;not null;uniqueIndex"`
	ConfigValue string `gorm:"type:text"`
}

// TableName 设置表名
func (SystemSetting) TableName() string {
	return "system_settings"
}

// SysSettings 系统设置管理结构
type SysSettings struct {
	db    *gorm.DB
	cache cache.Cache[string]
}

// NewSysSettings 创建系统设置管理器
func NewSysSettings(db *gorm.DB, redisClient *redis.Client) *SysSettings {
	// 使用Redis缓存存储系统设置，设置默认过期时间为7天，清理间隔为1小时
	redisCache := cache.NewRedisCache[string](redisClient, "sys:settings", time.Hour*24*7, time.Hour*1)

	return &SysSettings{
		db:    db,
		cache: redisCache,
	}
}

// 数据库操作方法
func (s *SysSettings) getFromDB(configName string) (string, error) {
	var setting SystemSetting
	result := s.db.Where("config_name = ?", configName).First(&setting)
	if result.Error != nil {
		return "", result.Error
	}
	return setting.ConfigValue, nil
}

func (s *SysSettings) setToDB(configName string, value string) error {
	var setting SystemSetting
	result := s.db.Where("config_name = ?", configName).First(&setting)

	if result.Error == nil {
		// 更新现有记录
		return s.db.Save(&setting).Error
	} else if result.Error == gorm.ErrRecordNotFound {
		// 创建新记录
		setting = SystemSetting{
			ConfigName:  configName,
			ConfigValue: value,
		}
		return s.db.Create(&setting).Error
	}

	return result.Error
}

// Get 获取字符串类型的配置值
func (s *SysSettings) Get(configName string) (string, error) {
	// 先从缓存获取
	if value, found := s.cache.Get(configName); found {
		return value, nil
	}

	// 缓存未命中，从数据库获取
	value, err := s.getFromDB(configName)
	if err != nil {
		return "", err
	}

	// 存入缓存
	s.cache.SetDefault(configName, value)
	return value, nil
}

// GetObj 获取结构体类型的配置值
func (s *SysSettings) GetObj(configName string, out interface{}) error {
	// 尝试从缓存获取
	if value, found := s.cache.Get(configName); found {
		return json.Unmarshal([]byte(value), out)
	}

	// 从数据库获取
	value, err := s.getFromDB(configName)
	if err != nil {
		return err
	}

	// 反序列化
	if err := json.Unmarshal([]byte(value), out); err != nil {
		return err
	}

	// 缓存
	s.cache.SetDefault(configName, value)
	return nil
}

// Set 设置配置值并更新缓存
func (s *SysSettings) Set(configName string, value interface{}) error {
	var valueStr string
	if strValue, ok := value.(string); ok {
		valueStr = strValue
	} else {
		bytes, err := json.Marshal(value)
		if err != nil {
			return err
		}
		valueStr = string(bytes)
	}

	// 更新数据库
	if err := s.setToDB(configName, valueStr); err != nil {
		return err
	}

	// 更新缓存
	s.cache.SetDefault(configName, valueStr)
	return nil
}

// Flush 清空缓存
func (s *SysSettings) Flush() {
	s.cache.Flush()
}

// RefreshCache 刷新指定配置的缓存
func (s *SysSettings) RefreshCache(configName string) error {
	value, err := s.getFromDB(configName)
	if err != nil {
		return err
	}

	s.cache.SetDefault(configName, value)
	return nil
}

// RefreshAllCache 刷新所有配置的缓存
func (s *SysSettings) RefreshAllCache() error {
	var settings []SystemSetting
	if err := s.db.Find(&settings).Error; err != nil {
		return err
	}

	s.cache.Flush()

	for _, setting := range settings {
		s.cache.SetDefault(setting.ConfigName, setting.ConfigValue)
	}

	return nil
}
