package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 初始化认证相关路由
func initAuthRouters(router *gin.Engine) {

	user := router.Group("/api/v1/user")
	{
		// 获取MFA二维码
		user.GET("/mfa_qrcode", controller.User.GetGoogleQrcode)
		// MFA认证
		user.POST("/mfa_auth", controller.User.GoogleQrcodeValidate)
		// 获取用户信息
		user.GET("/info", controller.User.GetUser)
	}

	// 登录
	router.POST("/api/auth/login", controller.User.Login)
	// 获取授权（钉钉）
	router.POST("/api/auth/dingtalk_login", controller.User.DingTalkLogin)
	// 获取授权（企业微信）
	router.POST("/api/auth/ww_login", controller.User.WeChatLogin)
	// 获取授权（飞书）
	router.POST("/api/auth/feishu_login", controller.User.FeishuLogin)
	// 注销
	router.POST("/api/auth/logout", controller.User.Logout)
}
