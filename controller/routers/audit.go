package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 初始化审计相关路由
func initAuditRouters(router *gin.Engine) {
	audit := router.Group("/api/v1/audit")
	{
		// 获取短信发送记录
		audit.GET("/sms", controller.Log.GetSMSRecordList)
		// 获取用户登录记录
		audit.GET("/login", controller.Login.GetLoginRecordList)
	}
}
