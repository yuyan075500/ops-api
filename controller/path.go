package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/service"
)

var Path path

type path struct{}

// GetPathListAll 获取所有接口
// @Summary 获取所有接口
// @Description 组相关接口
// @Tags 组管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/path/list [get]
func (p *path) GetPathListAll(c *gin.Context) {

	data, err := service.Path.GetPathListAll()

	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}

// GetPathList 获取接口列表（表格展示）
// @Summary 获取接口列表（表格展示）
// @Description 接口相关接口
// @Tags 接口管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param menu_name query string true "菜单名称"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/paths [get]
func (p *path) GetPathList(c *gin.Context) {
	params := new(struct {
		MenuName string `form:"menu_name" binding:"required"`
		Page     int    `form:"page" binding:"required"`
		Limit    int    `form:"limit" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	data, err := service.Path.GetPathList(params.MenuName, params.Page, params.Limit)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}
