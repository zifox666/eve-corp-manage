package database

import (
	"eve-corp-manager/config"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type SdeDbClient interface {
	Connect() (db *gorm.DB, err error)
}

type SQLiteConfig struct {
	Path string
}

// SdeDbInit 初始化SDE数据库连接
func SdeDbInit(dbClient SdeDbClient) (db *gorm.DB, dbErr error) {
	db, dbErr = dbClient.Connect()
	if dbErr != nil {
		return
	}
	return
}

// Connect SQLite连接
func (s *SQLiteConfig) Connect() (db *gorm.DB, err error) {
	dbPath := filepath.Join(config.DataDir, config.AppConfig.SdeSqlite.Path)

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: GetLogger(),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SQLite连接设置
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)

	return db, nil
}
