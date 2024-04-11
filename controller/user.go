package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/db"
	"ops-api/middleware"
	"ops-api/model"
	"ops-api/service"
	"time"
)

var User user

type user struct{}

// Login 用户登录
func (u *user) Login(c *gin.Context) {
	var (
		user model.AuthUser
		err  error
	)

	params := new(struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	})

	if err = c.Bind(params); err != nil {
		logger.Error("无效的请求参数：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 根据用户名查询用户
	if err := db.GORM.Where("username = ?", params.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 4404,
			"msg":  "用户不存在",
		})
		return
	}

	// 判断用户是否禁用
	if user.IsActive == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 4403,
			"msg":  "用户未激活，请联系管理员",
		})
		return
	}

	// 检查密码
	if user.CheckPassword(params.Password) == false {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 4404,
			"msg":  "用户密码错误",
		})
		return
	}

	token, err := middleware.GenerateJWT(user.ID, user.Name, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  "生成Token错误",
		})
		return
	}

	// 记录用户最后登录时间（待完成）
	db.GORM.Model(&user).Where("id = ?", user.ID).Update("last_login_at", time.Now())

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "认证成功",
		"token": token,
	})
}

// GetUserList 获取用户列表
func (u *user) GetUserList(c *gin.Context) {
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
		logger.Error("无效的请求参数：" + err.Error())
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
