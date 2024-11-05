package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/dao"
	"ops-api/service"
	"strconv"
)

var Account account

type account struct{}

// AddAccount 新增账号
// @Summary 新增账号
// @Description 账号相关接口
// @Tags 账号管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body service.AccountCreate true "账号信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建成功", "data": nil}"
// @Router /api/v1/account [post]
func (a *account) AddAccount(c *gin.Context) {
	var account = &service.AccountCreate{}

	if err := c.ShouldBind(account); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	userID := c.GetUint("id")
	if err := service.Account.AddAccount(account, userID); err != nil {
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

// DeleteAccount 删除账号
// @Summary 删除账号
// @Description 账号相关接口
// @Tags 账号管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "账号ID"
// @Success 200 {string} json "{"code": 0, "msg": "删除成功", "data": nil}"
// @Router /api/v1/account/{id} [delete]
func (a *account) DeleteAccount(c *gin.Context) {

	// 获取账号ID
	accountId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("ERROR：", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// 获取用户ID
	userID := c.GetUint("id")
	if err := service.Account.DeleteAccount(accountId, int(userID)); err != nil {
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

// UpdateAccount 更新账号信息
// @Summary 更新账号信息
// @Description 账号相关接口
// @Tags 账号管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param task body dao.AccountUpdate true "账号信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/account [put]
func (a *account) UpdateAccount(c *gin.Context) {
	var data = &dao.AccountUpdate{}

	if err := c.ShouldBind(&data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	userID := c.GetUint("id")
	if err := service.Account.UpdateAccount(data, userID); err != nil {
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

// GetAccountList 获取账号列表（自己的和别人分享的）
// @Summary 获取账号列表（自己的和别人分享的）
// @Description 账号相关接口
// @Tags 账号管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param name query string false "账号信息"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/accounts [get]
func (a *account) GetAccountList(c *gin.Context) {
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

	userID := c.GetUint("id")
	data, err := service.Account.GetAccountList(params.Name, userID, params.Page, params.Limit)
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

// GetAccountPassword 获取账号密码
// @Summary 获取账号密码
// @Description 账号相关接口
// @Tags 账号管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "账号ID"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/account/password/{id} [delete]
func (a *account) GetAccountPassword(c *gin.Context) {

	// 获取账号ID
	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("ERROR：", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	userId := c.GetUint("id")
	password, err := service.Account.GetAccountPassword(uint(accountID), username.(string), userId)
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
		"msg":  "获取成功，请妥善使用",
		"data": password,
	})
}

// GetSMSCode 获取验证码
// @Summary 获取验证码
// @Description 账号相关接口
// @Tags 账号管理
// @Accept application/json
// @Produce application/json
// @Param user body service.RestPassword true "用户信息"
// @Success 200 {string} json "{"code": 0, "msg": "校验码已发送，5分钟之内有效"}"
// @Router /api/v1/account/code [get]
func (a *account) GetSMSCode(c *gin.Context) {

	// 获取短信验证码
	userID := c.GetUint("id")
	if err := service.Account.GetSMSCode(userID); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "校验码已发送，5分钟之内有效",
	})
}

// CodeVerification 校验验证码
// @Summary 校验验证码
// @Description 账号相关接口
// @Tags 账号管理
// @Param user body service.CodeVerification true "验证码信息"
// @Success 200 {string} json "{"code": 0, "msg": "验证成功，本次验证有效期为10分钟"}"
// @Router /api/v1/account/code_verification [post]
func (a *account) CodeVerification(c *gin.Context) {
	var data = &service.CodeVerification{}

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
	userID := c.GetUint("id")
	if err := service.Account.CodeVerification(userID, data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "验证成功，本次验证有效期为10分钟",
	})
}
