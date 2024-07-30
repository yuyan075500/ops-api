package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 初始化菜单相关路由
func initMenuRouters(router *gin.Engine) {
	// 获取菜单列表
	router.GET("/api/v1/menus", controller.Menu.GetMenuList)
	// 获取接口列表
	router.GET("/api/v1/paths", controller.Path.GetPathList)
}
