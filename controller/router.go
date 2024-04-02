package controller

import (
	"github.com/gin-gonic/gin"
)

var Router router

type router struct{}

func (r *router) InitApiRouter(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.String(200, "项目初始化")
	})
}
