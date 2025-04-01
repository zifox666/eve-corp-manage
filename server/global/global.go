package global

import (
	"eve-corp-manager/core/system"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Logger          *zap.SugaredLogger
	LoggerLevel     = zap.NewAtomicLevel() // 支持通过http以及配置文件动态修改日志级别
	Db              *gorm.DB
	SdeDb           *gorm.DB
	Redis           *redis.Client
	Settings        *system.SysSettings
	Qq_notification bool
)
