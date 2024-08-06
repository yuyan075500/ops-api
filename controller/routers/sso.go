package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 单点登录相关路由
func initSSORouters(router *gin.Engine) {

	oauth := router.Group("/api/v1/oauth")
	{
		// 获取授权
		oauth.POST("/authorize", controller.SSO.Authorize)
		// 获取Token
		oauth.POST("/token", controller.SSO.GetToken)
		// 获取用户信息
		oauth.GET("/userinfo", controller.SSO.GetUserInfo)
	}
}
