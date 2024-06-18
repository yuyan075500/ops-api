package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
	"net/http"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/middleware"
	"ops-api/model"
	"ops-api/service"
	"ops-api/utils"
	"ops-api/utils/sms"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var User user

type user struct{}

// Login 登录
// @Summary 登录
// @Description 认证相关接口
// @Tags 用户认证
// @Accept application/json
// @Produce application/json
// @Param user body service.UserLogin true "用户名密码"
// @Success 200 {string} json "{"code": 0, "token": "用户令牌"}"
// @Router /login [post]
func (u *user) Login(c *gin.Context) {
	var (
		user   = &model.AuthUser{}
		params = &service.UserLogin{}
	)

	if err := c.ShouldBind(params); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	if params.LDAP {
		// 生成LDAP连接
		conn, err := middleware.CreateLDAPService()
		if err != nil {
			logger.Error("ERROR：" + err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": 90500,
				"msg":  err.Error(),
			})
			return
		}
		// 用户认证
		userInfo, err := conn.LDAPUserAuthentication(params.Username, params.Password)
		if err != nil {
			logger.Error("ERROR：" + err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": 90500,
				"msg":  err.Error(),
			})
			return
		}

		// 同步用户信息到本地数据库
		if err := global.MySQLClient.Where("username = ? AND user_from = ?", params.Username, "AD域").First(&user).Error; err != nil {
			// 如果登录类型为LDAP且用户不存在，则创建用户
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := service.User.AddUser(userInfo); err != nil {
					logger.Error("ERROR：" + err.Error())
					c.JSON(http.StatusOK, gin.H{
						"code": 90500,
						"msg":  err.Error(),
					})
					return
				}
			}
		}
	} else {
		// 根据用户名查询用户
		if err := global.MySQLClient.Where("username = ? AND user_from = ?", params.Username, "本地").First(&user).Error; err != nil {
			logger.Error("ERROR：" + err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": 90404,
				"msg":  "用户不存在",
			})
			return
		}
	}

	// 判断用户是否禁用
	if *user.IsActive == false {
		c.JSON(http.StatusOK, gin.H{
			"code": 90403,
			"msg":  "拒绝登录，请联系管理员",
		})
		return
	}

	// 检查密码
	if !params.LDAP {
		if user.CheckPassword(params.Password) == false {
			c.JSON(http.StatusOK, gin.H{
				"code": 90401,
				"msg":  "用户密码错误",
			})
			return
		}
	}

	token, err := middleware.GenerateJWT(user.ID, user.Name, user.Username)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// 记录用户最后登录时间
	global.MySQLClient.Model(&user).Where("id = ?", user.ID).Update("last_login_at", time.Now())

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"token": token,
	})
}

// Logout 注销
// @Summary 注销
// @Description 认证相关接口
// @Tags 用户认证
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json "{"code": 0, "data": nil}"
// @Router /logout [post]
func (u *user) Logout(c *gin.Context) {
	// 获取Token
	token := c.Request.Header.Get("Authorization")
	parts := strings.SplitN(token, " ", 2)

	// 将Token存入Redis缓存
	err := global.RedisClient.Set(parts[1], true, time.Duration(config.Conf.JWT.Expires)*time.Hour).Err()
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": nil,
	})
}

// UploadAvatar 头像上传
// @Summary 头像上传
// @Description 用户相关接口
// @Tags 用户管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param avatar formData file true "头像"
// @Success 200 {string} json "{"code": 0, "data": nil}"
// @Router /api/v1/user/avatarUpload [post]
func (u *user) UploadAvatar(c *gin.Context) {
	// 获取上传的头像
	avatar, err := c.FormFile("avatar")
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// 打开上传头像
	src, err := avatar.Open()
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
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
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// 将头像地址存储到数据库
	var user model.AuthUser
	global.MySQLClient.Model(&user).Where("username = ?", username).Update("avatar", avatarName)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": nil,
	})
}

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Description 认证相关接口
// @Tags 用户认证
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json "{"code": 0, "data": {}}"
// @Router /api/v1/user/info [get]
func (u *user) GetUser(c *gin.Context) {

	// 获取用户信息
	data, err := service.User.GetUser(c.GetUint("id"))
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// 返回用户信息
	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}

// GetUserListAll 获取所有的用户列表
// @Summary 获取所有的用户列表
// @Description 用户相关接口
// @Tags 用户管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/user/list [get]
func (u *user) GetUserListAll(c *gin.Context) {

	data, err := service.User.GetUserListAll()
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	err = sms.Send(
		"8822053031549",
		"13357110502",
		"770335531f48463283a478b179652f62",
		"[\"yuyan\", \"111111\"]",
		config.Conf.SMS.CallbackUrl,
		"物联亿达",
		"密码更改",
	)
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

// GetUserList 获取查询的用户列表
// @Summary 获取查询的用户列表
// @Description 用户相关接口
// @Tags 用户管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param name query string false "用户姓名"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/users [get]
func (u *user) GetUserList(c *gin.Context) {
	params := new(struct {
		Name  string `form:"name"`
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

	data, err := service.User.GetUserList(params.Name, params.Page, params.Limit)
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

// AddUser 创建用户
// @Summary 创建用户
// @Description 用户相关接口
// @Tags 用户管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body service.UserCreate true "用户信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建成功", "data": nil}"
// @Router /api/v1/user [post]
func (u *user) AddUser(c *gin.Context) {
	var user = &service.UserCreate{}

	if err := c.ShouldBind(user); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	if err := service.User.AddUser(user); err != nil {
		logger.Error("ERROR：" + err.Error())
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

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 用户相关接口
// @Tags 用户管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "用户ID"
// @Success 200 {string} json "{"code": 0, "msg": "删除成功", "data": nil}"
// @Router /api/v1/user/{id} [delete]
func (u *user) DeleteUser(c *gin.Context) {

	// 对ID进行类型转换
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("ERROR：", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// 执行删除
	if err := service.User.DeleteUser(userID); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
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

// UpdateUser 更新用户信息
// @Summary 更新用户信息
// @Description 用户相关接口
// @Tags 用户管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body dao.UserUpdate true "用户信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/user [put]
func (u *user) UpdateUser(c *gin.Context) {
	var data = &dao.UserUpdate{}

	// 解析请求参数
	if err := c.ShouldBind(&data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// 更新用户信息
	if err := service.User.UpdateUser(data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
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

// UpdateUserPassword 密码更新
// @Summary 密码更新
// @Description 用户相关接口
// @Tags 用户管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body dao.UserPasswordUpdate true "用户信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/user/reset_password [put]
func (u *user) UpdateUserPassword(c *gin.Context) {
	var data = &dao.UserPasswordUpdate{}

	// 解析请求参数
	if err := c.ShouldBind(&data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// 更新用户信息
	if err := service.User.UpdateUserPassword(data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
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

// ResetUserMFA MFA重置
// @Summary MFA重置
// @Description 用户相关接口
// @Tags 用户管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "用户ID"
// @Success 200 {string} json "{"code": 0, "msg": "重置成功", "data": nil}"
// @Router /api/v1/user/reset_mfa/{id} [put]
func (u *user) ResetUserMFA(c *gin.Context) {

	// 对ID进行类型转换
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// 更新用户信息
	if err := service.User.ResetUserMFA(userID); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "重置成功",
		"data": nil,
	})
}
