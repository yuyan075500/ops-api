package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/service"
	"strconv"
)

var Group group

type group struct{}

// GetGroupList 获取组列表
// @Summary 获取组列表
// @Description 组相关接口
// @Tags 组管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param name query string false "组名称"
// @Success 200 {string} json "{"code": 0, "msg": "获取列表成功", "data": []}"
// @Router /api/v1/groups [get]
func (u *group) GetGroupList(c *gin.Context) {
	params := new(struct {
		Name  string `form:"name"`
		Page  int    `form:"page" binding:"required"`
		Limit int    `form:"limit" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("无效的请求参数：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  "无效的请求参数",
		})
		return
	}

	data, err := service.Group.GetGroupList(params.Name, params.Page, params.Limit)
	if err != nil {
		logger.Error("获取列表失败：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "获取列表成功",
		"data": data,
	})
}

// AddGroup 创建组
// @Summary 创建组
// @Description 组相关接口
// @Tags 组管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param group body service.GroupCreate true "组信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建成功", "data": nil}"
// @Router /api/v1/group [post]
func (u *group) AddGroup(c *gin.Context) {
	var group = &service.GroupCreate{}

	if err := c.ShouldBind(group); err != nil {
		logger.Error("无效的请求参数：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	if err := service.Group.AddGroup(group); err != nil {
		logger.Error("新增分组失败：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "创建成功",
		"data": nil,
	})
}

// DeleteGroup 删除组
// @Summary 删除组
// @Description 组相关接口
// @Tags 组管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "组ID"
// @Success 200 {string} json "{"code": 0, "msg": "删除成功", "data": nil}"
// @Router /api/v1/group/{id} [delete]
func (u *group) DeleteGroup(c *gin.Context) {

	// 对ID进行类型转换
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("无效的分组ID：", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "无效的分组ID",
		})
		return
	}

	// 执行删除
	if err := service.Group.DeleteGroup(groupID); err != nil {
		logger.Error("删除失败：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "删除成功",
		"data": nil,
	})
}

// UpdateGroup 更新组信息
// @Summary 更新组信息
// @Description 组相关接口
// @Tags 组管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param group body service.GroupUpdate true "组信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/group [put]
func (u *group) UpdateGroup(c *gin.Context) {
	var data = &service.GroupUpdate{}

	// 解析请求参数
	if err := c.ShouldBind(&data); err != nil {
		logger.Error("无效的请求参数：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	// 更新用户信息
	if err := service.Group.UpdateGroup(data); err != nil {
		logger.Error("更新失败：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "更新成功",
		"data": nil,
	})
}

// UpdateGroupUser 更新组用户
// @Summary 更新组用户
// @Description 组相关接口
// @Tags 组管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param users body service.GroupUpdateUser true "用户信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/group/users [put]
func (u *group) UpdateGroupUser(c *gin.Context) {
	var data = &service.GroupUpdateUser{}

	// 解析请求参数
	if err := c.ShouldBind(&data); err != nil {
		logger.Error("无效的请求参数：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	// 更新用户信息
	if err := service.Group.UpdateGroupUser(data); err != nil {
		logger.Error("更新失败：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "更新成功",
		"data": nil,
	})
}
