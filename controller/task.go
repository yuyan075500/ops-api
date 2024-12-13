package controller

import (
	"github.com/gin-gonic/gin"
	"ops-api/dao"
	"ops-api/service"
	"strconv"
)

var Task task

type task struct{}

// AddTask 创建定时任务
// @Summary 创建定时任务
// @Description 定时任务相关接口
// @Tags 定时任务管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param task body service.TaskCreate true "定时任务信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建成功", "data": nil}"
// @Router /api/v1/task [post]
func (t *task) AddTask(c *gin.Context) {
	var task = &service.TaskCreate{}

	if err := c.ShouldBind(task); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	job, err := service.Task.AddTask(task)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "创建成功", job)
}

// DeleteTask 删除定时任务
// @Summary 删除定时任务
// @Description 定时任务相关接口
// @Tags 定时任务管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "定时任务ID"
// @Success 200 {string} json "{"code": 0, "msg": "删除成功"}"
// @Router /api/v1/task/{id} [delete]
func (t *task) DeleteTask(c *gin.Context) {

	// 对ID进行类型转换
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	// 执行删除
	if err := service.Task.DeleteTask(taskID); err != nil {
		Response(c, 90500, err.Error())
		return
	}

	Response(c, 0, "删除成功")
}

// UpdateTask 更新定时任务信息
// @Summary 更新定时任务信息
// @Description 定时任务相关接口
// @Tags 定时任务管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param task body dao.TaskUpdate true "定时任务信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/task [put]
func (t *task) UpdateTask(c *gin.Context) {
	var data = &dao.TaskUpdate{}

	// 解析请求参数
	if err := c.ShouldBind(&data); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	// 更新用户信息
	task, err := service.Task.UpdateTask(data)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	CreateOrUpdateResponse(c, 0, "更新成功", task)
}

// GetTaskList 获取定时任务列表
// @Summary 获取定时任务列表
// @Description 定时任务相关接口
// @Tags 定时任务管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param name query string false "定时任务名称"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/tasks [get]
func (t *task) GetTaskList(c *gin.Context) {
	params := new(struct {
		Name  string `form:"name"`
		Page  int    `form:"page" binding:"required"`
		Limit int    `form:"limit" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	data, err := service.Task.GetTaskList(params.Name, params.Page, params.Limit)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}

// GetTaskLogList 获取定时任务执行日志列表
// @Summary 获取定时任务执行日志列表
// @Description 定时任务相关接口
// @Tags 定时任务管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param id query int true "定时ID"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/task/logs [get]
func (t *task) GetTaskLogList(c *gin.Context) {
	params := new(struct {
		Id    uint `form:"id" binding:"required"`
		Page  int  `form:"page" binding:"required"`
		Limit int  `form:"limit" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		Response(c, 90400, err.Error())
		return
	}

	data, err := service.Task.GetTaskLogList(params.Id, params.Page, params.Limit)
	if err != nil {
		Response(c, 90500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}
