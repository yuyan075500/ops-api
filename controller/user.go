package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/global"
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

// Login 登录
// @Summary 登录
// @Description 用户认证相关接口
// @Tags 用户认证
// @Accept application/json
// @Produce application/json
// @Param user body service.UserLogin true "用户名密码"
// @Success 200 {string} json "{"code": 0, "token": "用户令牌"}"
// @Router /login [post]
func (u *user) Login(c *gin.Context) {
	var params = &service.UserLogin{}

	if err := c.ShouldBind(params); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	token, redirectUri, redirect, err := service.User.Login(params, c)
	if err != nil {

		// 开启事务
		tx := global.MySQLClient.Begin()

		// 新增登录失败记录
		if err := service.Login.AddLoginRecord(tx, 2, params.Username, "账号密码", err, c); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 90500,
				"msg":  err.Error(),
			})
			return
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return
		}

		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// 如果开启MFA认证需要返回重定向的页面
	if redirect != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":     0,
			"token":    token,
			"redirect": redirect,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":         0,
		"token":        token,
		"redirect_uri": redirectUri,
	})
}

// Logout 注销
// @Summary 注销
// @Description 用户认证相关接口
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
// @Description 个人信息管理相关接口
// @Tags 个人信息管理
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
		"msg":  "头像更新成功",
		"data": nil,
	})
}

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Description 用户认证相关接口
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

// GetUserListAll 获取用户列表（下拉框）
// @Summary 获取用户列表（下拉框）
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
// @Param user body dao.UserCreate true "用户信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建成功", "data": nil}"
// @Router /api/v1/user [post]
func (u *user) AddUser(c *gin.Context) {
	var user = &dao.UserCreate{}

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

// GetVerificationCode 获取验证码
// @Summary 获取验证码
// @Description 个人信息管理相关接口
// @Tags 个人信息管理
// @Accept application/json
// @Produce application/json
// @Param user body service.RestPassword true "用户信息"
// @Success 200 {string} json "{"code": 0, "msg": "校验码已发送..."}"
// @Router /api/v1/sms/reset_password_code [post]
func (u *user) GetVerificationCode(c *gin.Context) {

	var (
		data           = &service.UserInfo{}
		expirationTime = 5 // 指定验证码过期时间
	)

	// 解析请求参数
	if err := c.ShouldBind(&data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// 获取短信验证码
	if err := service.User.GetVerificationCode(data, expirationTime); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  fmt.Sprintf("校验码已发送，%s分钟之内有效", strconv.Itoa(expirationTime)),
	})
}

// UpdateSelfPassword 密码更新
// @Summary 密码更新
// @Description 个人信息管理相关接口
// @Tags 个人信息管理
// @Param user body dao.UserPasswordUpdate true "用户信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功"}"
// @Router /api/v1/reset_password [post]
func (u *user) UpdateSelfPassword(c *gin.Context) {
	var data = &service.RestPassword{}

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
	if err := service.User.UpdateSelfPassword(data); err != nil {
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
	})
}

// UserSyncAd AD域用户同步
// @Summary AD域用户同步
// @Description 用户相关接口
// @Tags 用户管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json "{"code": 0, "msg": "同步成功"}"
// @Router /api/v1/user/sync/ad [post]
func (u *user) UserSyncAd(c *gin.Context) {

	// 同步用户
	if err := service.User.UserSync(); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "同步成功",
	})
}

// GetGoogleQrcode 获取MFA二维码
// @Summary 获取MFA二维码
// @Description 用户认证相关接口
// @Tags 用户认证
// @Param token query string true "用户认证通过后的Token"
// @Success 200 {string} json "{"code": 0, "qrcode": ""}"
// @Router /api/v1/user/mfa_qrcode [get]
func (u *user) GetGoogleQrcode(c *gin.Context) {
	params := new(struct {
		Token string `form:"token"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// 获取二维码
	qrcode, err := service.MFA.GetGoogleQrcode(params.Token)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// 返回用户信息
	c.Header("Content-Type", "image/png")
	c.JSON(200, gin.H{
		"code":   0,
		"qrcode": qrcode,
	})
}

// GoogleQrcodeValidate MFA认证
// @Summary MFA认证
// @Description 用户认证相关接口
// @Tags 用户认证
// @Param user body service.MFAValidate true "MFA认证信息"
// @Success 200 {string} json "{"code": 0, "token": "用户令牌"}"
// @Router /api/v1/user/mfa_auth [post]
func (u *user) GoogleQrcodeValidate(c *gin.Context) {

	var params = &service.MFAValidate{}

	// 请求参数绑定
	if err := c.ShouldBind(params); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// MFA校验
	token, err := service.MFA.GoogleQrcodeValidate(params.Username, params.Token, params.Code, c)
	if err != nil {

		// 开启事务
		tx := global.MySQLClient.Begin()

		// 新增登录记录
		if err := service.Login.AddLoginRecord(tx, 2, params.Username, "双因子", err, c); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 90500,
				"msg":  err.Error(),
			})
			return
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return
		}

		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":  0,
		"token": token,
	})
}
