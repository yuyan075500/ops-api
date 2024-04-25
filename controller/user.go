package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/config"
	"ops-api/global"
	"ops-api/middleware"
	"ops-api/model"
	"ops-api/service"
	"ops-api/utils"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var User user

type user struct{}

// Login 用户登录
// @Summary 用户登录
// @Description 用户相关接口
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Param user body service.UserLogin true "用户名密码"
// @Success 200 {string} json "{"code": 0, "msg": "认证成功", "token": "用户令牌"}"
// @Router /login [post]
func (u *user) Login(c *gin.Context) {
	var (
		user   = &model.AuthUser{}
		params = &service.UserLogin{}
	)

	if err := c.ShouldBind(params); err != nil {
		logger.Error("无效的请求参数：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 根据用户名查询用户
	if err := global.MySQLClient.Where("username = ?", params.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 4404,
			"msg":  "用户不存在",
		})
		return
	}

	// 判断用户是否禁用
	if *user.IsActive == false {
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

	// 记录用户最后登录时间
	global.MySQLClient.Model(&user).Where("id = ?", user.ID).Update("last_login_at", time.Now())

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "认证成功",
		"token": token,
	})
}

// Logout 用户注销
// @Summary 用户注销
// @Description 用户相关接口
// @Tags 用户管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json "{"code": 0, "msg": "注销成功", "data": nil}"
// @Router /logout [post]
func (u *user) Logout(c *gin.Context) {
	// 获取Token
	token := c.Request.Header.Get("Authorization")
	parts := strings.SplitN(token, " ", 2)

	// 将Token存入Redis缓存
	err := global.RedisClient.Set(parts[1], true, time.Duration(config.Conf.JWT.Expires)*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  "用户注销失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "注销成功",
		"data": nil,
	})
}

// UploadAvatar 用户头像上传
// @Summary 用户头像上传
// @Description 用户相关接口
// @Tags 用户管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param avatar formData file true "头像"
// @Success 200 {string} json "{"code": 0, "msg": "头像更新成功", "data": nil}"
// @Router /api/v1/user/avatarUpload [post]
func (u *user) UploadAvatar(c *gin.Context) {
	// 获取上传的头像
	avatar, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4000,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 打开上传头像
	src, err := avatar.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4000,
			"msg":  "文件打开失败",
		})
		return
	}

	// 上传头像
	// 获取当前登录用户的用户名
	username, _ := c.Get("username")
	// 拼接头像存储的路径和文件名：avatar/<用户名><文件后缀>
	avatarName := fmt.Sprintf("avatar/%v%v", username, filepath.Ext(avatar.Filename))
	err = utils.FileUpload(avatarName, avatar.Header.Get("Content-Type"), src, avatar.Size)
	if err != nil {
		logger.Error("文件上传失败：" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4000,
			"msg":  "文件上传失败",
		})
		return
	}

	// 将头像地址存储到数据库
	var user model.AuthUser
	global.MySQLClient.Model(&user).Where("username = ?", username).Update("avatar", avatarName)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "头像更新成功",
		"data": nil,
	})
}

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Description 用户相关接口
// @Tags 用户管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json "{"code": 0, "msg": "获取用户信息成功", "data": {}}"
// @Router /api/v1/user/info [get]
func (u *user) GetUser(c *gin.Context) {

	// 获取用户信息
	data, err := service.User.GetUser(c.GetUint("id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 90404,
			"msg":  "获取用户信息失败",
		})
		return
	}

	// 返回用户信息
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "获取用户信息成功",
		"data": data,
	})
}

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 用户相关接口
// @Tags 用户管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param name query string false "用户姓名"
// @Success 200 {string} json "{"code": 0, "msg": "获取用户列表成功", "data": []}"
// @Router /api/v1/users [get]
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
		"msg":  "获取用户列表成功",
		"data": data,
	})
}

// AddUser 创建用户
// @Summary 创建用户
// @Description 用户相关接口
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body service.UserCreate true "用户信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建用户成功", "data": nil}"
// @Router /api/v1/user [post]
func (u *user) AddUser(c *gin.Context) {
	var user = &service.UserCreate{}

	if err := c.ShouldBind(user); err != nil {
		logger.Error("无效的请求参数：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	if err := service.User.AddUser(user); err != nil {
		logger.Error("新增用户失败：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "创建用户成功",
		"data": nil,
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 用户相关接口
// @Tags 用户管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "用户ID"
// @Success 200 {string} json "{"code": 0, "msg": "删除用户成功", "data": nil}"
// @Router /api/v1/user/{id} [delete]
func (u *user) DeleteUser(c *gin.Context) {

	// 对ID进行类型转换
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("无效的用户ID：", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 4001,
			"msg":  "无效的用户ID",
		})
		return
	}

	// 执行删除
	if err := service.User.DeleteUser(userID); err != nil {
		logger.Error("删除用户失败：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "删除用户成功",
		"data": nil,
	})
}

// UpdateUser 用户更新信息
// @Summary 用户更新信息
// @Description 用户相关接口
// @Tags 用户管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body service.UserUpdate true "用户信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新用户成功", "data": nil}"
// @Router /api/v1/user [put]
func (u *user) UpdateUser(c *gin.Context) {
	var data = &service.UserUpdate{}

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
	if err := service.User.UpdateUser(data); err != nil {
		logger.Error("更新用户失败：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "更新用户成功",
		"data": nil,
	})
}
