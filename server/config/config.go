package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	App struct {
		Env  string
		Port string
	}
	Database struct {
		Dsn          string
		MaxIdleConns int
		MaxOpenConns int
		WaitTimeOut  int
	}
	Redis struct {
		Addr     string
		DB       int
		Password string
	}
	SourceDir struct {
		Base string
	}
	SdeSqlite struct {
		Path string
	}
	Proxy struct {
		Host string
		Port string
	}
	QQ struct {
		OnebotUrl string
	}
}

var AppConfig *Config

var (
	RootDir string
	LogsDir string
	DataDir string
)

func InitRootDir() {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	RootDir = filepath.Dir(rootDir)
	DataDir = filepath.Join(RootDir, "data")
	LogsDir = filepath.Join(DataDir, "logs")
}

func InitConfig() {
	InitRootDir()

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join(RootDir, "server"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	AppConfig = &Config{}

	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}
