package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 单点登录相关路由
func initSSORouters(router *gin.Engine) {

	sso := router.Group("/api/v1/sso")
	{
		// 获取授权（OAuth2.0）
		sso.POST("/oauth/authorize", controller.SSO.OAuthAuthorize)
		// 获取授权（CAS3.0）
		sso.POST("/cas/authorize", controller.SSO.CASAuthorize)
		// 获取Token（OAuth2.0）
		sso.POST("/token", controller.SSO.GetToken)
		// 获取用户信息（OAuth2.0）
		sso.GET("/userinfo", controller.SSO.GetUserInfo)
	}

	// CAS3.0客户端票据校验
	router.GET("/p3/serviceValidate", controller.SSO.CASServiceValidate)
}
