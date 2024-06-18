package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"io"
	"ops-api/service"
)

var Log log

type log struct{}

func (l *log) SMSCallback(c *gin.Context) {

	// 获取回调请求Body中的内容
	body, _ := io.ReadAll(c.Request.Body)
	bodyStr := fmt.Sprintf("%s", body)

	if err := service.Log.SMSCallback(bodyStr); err != nil {
		logger.Error("ERROR：" + err.Error())
		return
	}
}
