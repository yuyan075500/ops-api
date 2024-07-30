package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 初始化短信相关路由
func initSmsRouters(router *gin.Engine) {
	sms := router.Group("/api/v1/sms")
	{
		// 获取重置密码验证码
		sms.POST("/reset_password_code", controller.User.GetVerificationCode)
		// 接收短信回调地址
		sms.POST("/callback", controller.Log.SMSCallback)
	}
}
