package corp_pap

import (
	"eve-corp-manager/api/v1/service"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(routerGroup *gin.RouterGroup) {
	// 创建corpPap路由组
	corpPapRouter := routerGroup.Group("corp_pap")
	{
		// 获取用户PAP记录列表
		corpPapRouter.GET("/list", service.GetUserPapList)
		// 获取用户PAP余额
		corpPapRouter.GET("/balance", service.GetUserPapBalance)
		// 增加用户PAP
		corpPapRouter.POST("/add", service.AddUserPap)
		// 消费用户PAP
		corpPapRouter.POST("/consume", service.ConsumeUserPap)
		// 获取PAP操作日志
		corpPapRouter.GET("/logs", service.GetPapLogs)
	}
}
