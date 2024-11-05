package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 初始化账号相关路由
func initAccountRouters(router *gin.Engine) {
	// 获取账号列表（表格）
	router.GET("/api/v1/accounts", controller.Account.GetAccountList)

	account := router.Group("/api/v1/account")
	{
		// 新增账号
		account.POST("", controller.Account.AddAccount)
		// 删除定时任务
		account.DELETE("/:id", controller.Account.DeleteAccount)
		// 修改定时任务
		account.PUT("", controller.Account.UpdateAccount)
		// 获取账号密码
		account.GET("/password/:id", controller.Account.GetAccountPassword)
		// 获取短信验证码
		account.GET("/code", controller.Account.GetSMSCode)
		// 获取短信验证码
		account.POST("/code_verification", controller.Account.CodeVerification)
	}
}
