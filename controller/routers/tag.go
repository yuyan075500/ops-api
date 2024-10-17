package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 初始化标签相关路由
func initTagRouters(router *gin.Engine) {
	tag := router.Group("/api/v1/tag")
	{
		// 获取标签列表
		tag.GET("/list", controller.Tag.GetTagList)
	}
}
