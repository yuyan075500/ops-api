package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 初始化短信相关路由
func initSmsRouters(router *gin.Engine) {
	sms := router.Group("/api/v1/sms")
	{
		// 接收短信回调（华为云）
		sms.POST("/huawei/callback", controller.SMS.SMSCallback)
	}
}
