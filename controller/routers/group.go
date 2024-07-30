package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 初始化用户分组相关路由
func initGroupRouters(router *gin.Engine) {
	// 获取分组列表
	router.GET("/api/v1/groups", controller.Group.GetGroupList)
	// 获取接口列表
	router.GET("/api/v1/path/list", controller.Path.GetPathListAll)
	// 获取菜单列表
	router.GET("/api/v1/menu/list", controller.Menu.GetMenuListAll)

	group := router.Group("/api/v1/group")
	{
		// 新增分组
		group.POST("", controller.Group.AddGroup)
		// 修改分组
		group.PUT("", controller.Group.UpdateGroup)
		// 修改分组用户
		group.PUT("/users", controller.Group.UpdateGroupUser)
		// 修改分组权限
		group.PUT("/permissions", controller.Group.UpdateGroupPermission)
		// 删除分组
		group.DELETE("/:id", controller.Group.DeleteGroup)
	}
}
