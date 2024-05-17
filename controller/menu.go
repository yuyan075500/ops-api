package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/service"
)

var Menu menu

type menu struct{}

// GetMenuListAll 获取所有的菜单列表
// @Summary 获取所有的菜单列表
// @Description 菜单关接口
// @Tags 菜单管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/menu/list [get]
func (u *menu) GetMenuListAll(c *gin.Context) {

	data, err := service.Menu.GetMenuListAll()

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

// GetMenuList 获取菜单列表
// @Summary 获取菜单列表
// @Description 菜单关接口
// @Tags 菜单管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/menus [get]
func (u *menu) GetMenuList(c *gin.Context) {
	params := new(struct {
		Title string `form:"title"`
		Page  int    `form:"page" binding:"required"`
		Limit int    `form:"limit" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	data, err := service.Menu.GetMenuList(params.Title, params.Page, params.Limit)
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

// GetUserMenu 获取用户菜单
// @Summary 获取用户菜单
// @Description 菜单关接口
// @Tags 菜单管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/user/menu [get]
func (u *menu) GetUserMenu(c *gin.Context) {

	data, err := service.Menu.GetUserMenu()
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
