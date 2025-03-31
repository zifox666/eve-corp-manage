package run_log

import (
	"eve-corp-manager/config"
	"eve-corp-manager/core/common"
	"eve-corp-manager/global"
	"os"

	"go.uber.org/zap"
)

func InitLog(runmode string, filePath string) (*zap.SugaredLogger, error) {

	runtimePath := config.RootDir + "/logs"
	if err := os.MkdirAll(runtimePath, 0777); err != nil {
		return nil, err
	}
	var level zap.AtomicLevel
	if runmode == "debug" {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		level = global.LoggerLevel
	}

	logger := common.InitLogger(runtimePath+filePath, level)
	return logger, nil
}
