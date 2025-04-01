package qq

import (
	"eve-corp-manager/config"
	"eve-corp-manager/core/qq"
	"eve-corp-manager/global"
)

// InitQQClient 初始化QQ客户端
func InitQQClient() {
	baseURL := config.AppConfig.QQ.OnebotUrl
	if baseURL == "" {
		global.Logger.Info("QQ通知服务地址未配置，不启用通知服务")
		global.Qq_notification = false
	} else {
		qq.QQClient = qq.NewClient(baseURL)
		global.Logger.Info("QQ通知服务初始化成功")
		global.Qq_notification = true
	}
}
