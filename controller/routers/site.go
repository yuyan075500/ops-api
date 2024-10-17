package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 初始化站点相关路由
func initSiteRouters(router *gin.Engine) {
	// 获取站点列表（表格）
	router.GET("/api/v1/sites", controller.Site.GetSiteList)

	site := router.Group("/api/v1/site")
	{
		// 新增站点
		site.POST("", controller.Site.AddSite)
		// 修改站点
		site.PUT("", controller.Site.UpdateSite)
		// 删除站点
		site.DELETE("/:id", controller.Site.DeleteSite)
		// 上传站点Logo
		site.POST("/logoUpload", controller.Site.UploadLogo)
		// 获取站点列表（导航页）
		site.GET("/guide", controller.Site.GetSiteGuideList)
		// 创建站点分组
		site.POST("/group", controller.Site.AddGroup)
		// 修改站点分组
		site.PUT("/group", controller.Site.UpdateGroup)
		// 删除站点分组
		site.DELETE("/group/:id", controller.Site.DeleteGroup)
		// 修改站点用户
		site.PUT("/users", controller.Site.UpdateSiteUser)
		// 修改站点标准
		site.PUT("/tags", controller.Site.UpdateSiteTag)
	}
}
