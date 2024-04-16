package controller

import (
	"github.com/gin-gonic/gin"
)

var Router router

type router struct{}

func (r *router) InitApiRouter(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Ok")
	}).
		POST("/login", User.Login).
		POST("/logout", User.Logout).
		GET("/api/v1/users", User.GetUserList).
		GET("/api/v1/user/info", User.GetUser).
		POST("/api/v1/user", User.AddUser).
		POST("/api/v1/user/avatarUpload", User.UploadAvatar)
}
