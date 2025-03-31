package initialize

import (
	"eve-corp-manager/config"
	"eve-corp-manager/global"
	"eve-corp-manager/initialize/run_log"
	"github.com/gin-gonic/gin"
	"log"
)

func StartUp() {
	gin.SetMode(config.AppConfig.App.Env)

	if logger, err := run_log.InitLog(config.AppConfig.App.Env, "/running_"+config.AppConfig.App.Env+".log"); err != nil {
		log.Panicln("Log initialization error", err)
	} else {
		global.Logger = logger
	}
}
