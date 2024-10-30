package routers

import (
	"github.com/gin-gonic/gin"
	"ops-api/controller"
)

// 初始定时任务相关路由
func initTaskRouters(router *gin.Engine) {
	// 获取定时任务列表（表格）
	router.GET("/api/v1/tasks", controller.Task.GetTaskList)

	task := router.Group("/api/v1/task")
	{
		// 新增定时任务
		task.POST("", controller.Task.AddTask)
		// 修改定时任务
		task.PUT("", controller.Task.UpdateTask)
		// 删除定时任务
		task.DELETE("/:id", controller.Task.DeleteTask)
		// 获取定时任务执行日志列表（表格）
		task.GET("/logs", controller.Task.GetTaskLogList)
	}
}
