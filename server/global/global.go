package global

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Logger      *zap.SugaredLogger
	LoggerLevel = zap.NewAtomicLevel() // 支持通过http以及配置文件动态修改日志级别
	Db          *gorm.DB
	Redis       *redis.Client
)
