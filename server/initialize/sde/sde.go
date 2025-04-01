package sde

import (
	"eve-corp-manager/config"
	"eve-corp-manager/core/sde"
	"eve-corp-manager/global"
	"fmt"
	"os"
	"path/filepath"
)

// InitSDE 初始化SDE数据库
func InitSDE() error {
	// 获取SDE SQLite文件路径
	sdePath := filepath.Join(config.DataDir, config.AppConfig.SdeSqlite.Path)

	// 检查文件是否存在
	_, err := os.Stat(sdePath)
	if os.IsNotExist(err) {
		global.Logger.Info("SDE数据库文件不存在，开始下载...")
		// 下载SDE数据库
		err = sde.DownloadSDE(sdePath)
		if err != nil {
			return fmt.Errorf("下载SDE数据库失败: %v", err)
		}
		global.Logger.Info("SDE数据库下载完成")
	}

	return nil
}
