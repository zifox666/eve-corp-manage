package initialize

import (
	"eve-corp-manager/config"
	"eve-corp-manager/global"
	"eve-corp-manager/initialize/database"
	"eve-corp-manager/initialize/redis"
	"eve-corp-manager/initialize/run_log"
	"eve-corp-manager/initialize/sde"
	"eve-corp-manager/initialize/system"
	"eve-corp-manager/models"
	"github.com/gin-gonic/gin"
	"log"
)

func StartUp() {
	gin.SetMode(config.AppConfig.App.Env)

	// 启动日志服务
	if logger, err := run_log.InitLog(config.AppConfig.App.Env, "/running_"+config.AppConfig.App.Env+".log"); err != nil {
		log.Panicln("Log initialization error", err)
	} else {
		global.Logger = logger
	}

	// 初始化 SDE 数据库
	if err := sde.InitSDE(); err != nil {
		log.Panicln("SDE 数据库初始化错误", err)
	}

	// 启动数据库服务
	startDb()

	// 启动Redis服务
	rdb, err := redis.InitRedis(redis.Options{
		Addr:     config.AppConfig.Redis.Addr,
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})
	if err != nil {
		log.Panicln("Redis initialization error", err)
	}
	global.Redis = rdb

	// 加载系统缓存
	system.InitSettings()

}

func startDb() {
	// 连接主数据库
	var dbClientInfo database.DbClient
	dbClientInfo = &database.MySQLConfig{
		Dsn:          config.AppConfig.Database.Dsn,
		MaxIdleConns: config.AppConfig.Database.MaxIdleConns,
		MaxOpenConns: config.AppConfig.Database.MaxOpenConns,
		WaitTimeout:  config.AppConfig.Database.WaitTimeOut,
	}
	if db, err := database.DbInit(dbClientInfo); err != nil {
		log.Panicln("Database initialization error", err)
	} else {
		global.Db = db
		models.Db = global.Db
	}
	err := database.CreateDatabase(global.Db)
	if err != nil {
		log.Panicln("Database migration error", err)
	}

	// 连接 SDE 数据库
	var sdeClientInfo database.SdeDbClient
	sdeClientInfo = &database.SQLiteConfig{
		Path: config.AppConfig.SdeSqlite.Path,
	}
	if sdeDb, err := database.SdeDbInit(sdeClientInfo); err != nil {
		log.Panicln("SDE 数据库连接错误", err)
	} else {
		global.SdeDb = sdeDb
	}
}
