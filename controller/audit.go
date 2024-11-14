package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/service"
)

var Audit audit

type audit struct{}

// GetSMSRecord 获取短信发送记录
// @Summary 获取短信发送记录
// @Description 审计相关接口
// @Tags 短信发送记录管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param receiver query string false "电话号码"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/audit/sms [get]
func (l *audit) GetSMSRecord(c *gin.Context) {

	params := new(struct {
		Receiver string `form:"receiver"`
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

	data, err := service.Audit.GetSMSRecordList(params.Receiver, params.Page, params.Limit)

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

// GetSMSReceipt 获取短信回执（阿里云）
// @Summary 获取短信回执（阿里云）
// @Description 审计相关接口
// @Tags 登录日志管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id query int true "短信记录ID"
// @Success 200 {string} json "{"code": 0}"
// @Router /api/v1/audit/sms/receipt [get]
func (l *audit) GetSMSReceipt(c *gin.Context) {
	params := new(struct {
		Id int `form:"id" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	err := service.Audit.GetSMSReceipt(params.Id)
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
	})
}

// GetLoginRecord 获取系统登录记录
// @Summary 获取系统登录记录
// @Description 审计相关接口
// @Tags 登录日志管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param name query string false "关键字"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/audit/login [get]
func (l *audit) GetLoginRecord(c *gin.Context) {

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

	data, err := service.Audit.GetLoginRecordList(params.Name, params.Page, params.Limit)

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

// GetOplog 获取系统操作记录
// @Summary 获取系统操作记录
// @Description 审计相关接口
// @Tags 操作日志管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param name query string false "关键字"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/audit/oplog [get]
func (l *audit) GetOplog(c *gin.Context) {

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

	data, err := service.Audit.GetOplogList(params.Name, params.Page, params.Limit)

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
