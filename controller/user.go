package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/service"
)

var User user

type user struct{}

// GetUserList 获取用户列表
func (u *user) GetUserList(c *gin.Context) {
	params := new(struct {
		Name  string `form:"name"`
		Page  int    `form:"page"`
		Limit int    `form:"limit"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("请求参数无效：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	data, err := service.User.GetUserList(params.Name, params.Page, params.Limit)
	if err != nil {
		logger.Error("获取用户列表失败：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}

// AddUser 创建用户
func (u *user) AddUser(c *gin.Context) {
	var (
		user = &service.UserCreate{}
		err  error
	)

	if err = c.ShouldBind(user); err != nil {
		logger.Error("无效的参数：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	if err = service.User.AddUser(user); err != nil {
		logger.Error("创建用户失败：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "创建用户成功",
	})
}
