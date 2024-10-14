package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/service"
)

var Login login

type login struct{}

// GetLoginRecordList 获取用户登录列表
// @Summary 获取用户登录列表
// @Description 日志相关接口
// @Tags 日志管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param username query string false "用户名"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/audit/login [get]
func (l *login) GetLoginRecordList(c *gin.Context) {

	params := new(struct {
		Username string `form:"username"`
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

	data, err := service.Login.GetLoginRecordList(params.Username, params.Page, params.Limit)

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

// GetSMSReceipt 获取短信回执
// @Summary 获取短信回执
// @Description 日志相关接口
// @Tags 日志管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id query int true "短信记录ID"
// @Success 200 {string} json "{"code": 0}"
// @Router /api/v1/audit/sms/receipt [get]
func (l *login) GetSMSReceipt(c *gin.Context) {
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

	err := service.Login.GetSMSReceipt(params.Id)
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
