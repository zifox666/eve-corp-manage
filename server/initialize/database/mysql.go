package database

import (
	"eve-corp-manager/config"
	system2 "eve-corp-manager/core/system"
	"eve-corp-manager/models/service"
	"eve-corp-manager/models/system"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DbClient interface {
	Connect() (db *gorm.DB, err error)
}

type MySQLConfig struct {
	Dsn          string
	MaxIdleConns int
	MaxOpenConns int
	WaitTimeout  int
}

func DbInit(dbClient DbClient) (db *gorm.DB, dbErr error) {
	db, dbErr = dbClient.Connect()
	if dbErr != nil {
		return
	}
	return
}

// Connect Mysql连接
func (d *MySQLConfig) Connect() (db *gorm.DB, err error) {
	dsn := config.AppConfig.Database.Dsn
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: GetLogger(),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(config.AppConfig.Database.MaxIdleConns)                                     // SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxOpenConns(config.AppConfig.Database.MaxOpenConns)                                     // SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDb.SetConnMaxLifetime(time.Duration(config.AppConfig.Database.WaitTimeOut * int(time.Second))) // SetConnMaxLifetime 设置了连接可复用的最大时间
	return
}

// 日志
func GetLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Warn, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 彩色打印
		},
	)

}

// 创建数据库
func CreateDatabase(db *gorm.DB) error {

	db = db.Set("gorm:table_options", "ENGINE=InnoDB")

	err := db.AutoMigrate(
		&system.User{},
		&system.Role{},
		&system.RoleMenu{},
		&system2.SystemSetting{},

		&service.Fleet{},
		&service.CharacterFleetAssociation{},
	)

	// 创建数据表
	err = db.AutoMigrate()

	return err
}
