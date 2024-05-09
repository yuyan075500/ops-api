package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/service"
)

var Menu menu

type menu struct{}

// GetUserMenu 获取用户菜单
// @Summary 获取用户菜单
// @Description 菜单关接口
// @Tags 菜单管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json "{"code": 0, "msg": "获取列表成功", "data": []}"
// @Router /api/v1/menus [get]
func (u *menu) GetUserMenu(c *gin.Context) {

	data, err := service.Menu.GetUserMenu()
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
