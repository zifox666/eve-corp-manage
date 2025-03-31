package system

import (
	"eve-corp-manager/core/system"
	"eve-corp-manager/global"
)

// InitSettings 初始化系统设置
func InitSettings() {
	global.Settings = system.NewSysSettings(global.Db, global.Redis)
}
