package controller

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"ops-api/config"
	"ops-api/docs"
)

var Router router

type router struct{}

func (r *router) InitRouter(router *gin.Engine) {

	// Swagger 接口
	if config.Conf.Swagger {
		docs.SwaggerInfo.BasePath = ""
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// 系统接口
	router.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	}).
		POST("/login", User.Login).
		POST("/logout", User.Logout).
		POST("/api/v1/user", User.AddUser).
		GET("/api/v1/user/info", User.GetUser).
		DELETE("/api/v1/user/:id", User.DeleteUser).
		GET("/api/v1/users", User.GetUserList).
		POST("/api/v1/user/avatarUpload", User.UploadAvatar)
}
