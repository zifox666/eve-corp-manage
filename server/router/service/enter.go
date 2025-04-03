package service

import (
	"eve-corp-manager/router/service/corp_pap"

	"github.com/gin-gonic/gin"
)

// Init 初始化服务模块路由
func Init(routerGroup *gin.RouterGroup) {
	// 服务路由组
	serviceRouter := routerGroup.Group("service")

	// 初始化各个服务模块的路由
	corp_pap.Init(serviceRouter)
	// 这里可以添加其他服务模块的路由初始化
}
