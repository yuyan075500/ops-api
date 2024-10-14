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
		audit.GET("/sms", controller.Audit.GetSMSRecord)
		// 获取短信回执（阿里云）
		audit.GET("/sms/receipt", controller.Audit.GetSMSReceipt)
		// 获取系统登录记录
		audit.GET("/login", controller.Audit.GetLoginRecord)
	}
}
