package router

import (
	"github.com/gin-gonic/gin"

	"eve-corp-manager/global"
	"eve-corp-manager/router/service"
	"eve-corp-manager/router/system"
)

func InitRouters(addr string) error {
	router := gin.Default()
	rootRouter := router.Group("/")
	routerGroup := rootRouter.Group("api/v1")

	// 接口
	system.Init(routerGroup)
	// 初始化服务模块路由
	service.Init(routerGroup)

	global.Logger.Info("eve-corp-manager 后端服务已经启动，监听 ", addr)
	return router.Run(addr)
}
