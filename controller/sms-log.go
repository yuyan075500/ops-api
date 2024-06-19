package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"io"
	"net/http"
	"ops-api/service"
)

var Log log

type log struct{}

// GetSMSRecordList 获取短信发送列表
// @Summary 获取短信发送列表
// @Description 日志相关接口
// @Tags 日志管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param receiver query string false "电话号码"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/audit/sms/list [get]
func (l *log) GetSMSRecordList(c *gin.Context) {

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

	data, err := service.Log.GetSMSRecordList(params.Receiver, params.Page, params.Limit)

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

// SMSCallback 短信回调
func (l *log) SMSCallback(c *gin.Context) {

	// 获取回调请求Body中的内容
	body, _ := io.ReadAll(c.Request.Body)
	bodyStr := fmt.Sprintf("%s", body)

	if err := service.Log.SMSCallback(bodyStr); err != nil {
		logger.Error("ERROR：" + err.Error())
		return
	}
}
