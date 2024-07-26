package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"ops-api/dao"
	"ops-api/model"
)

var Login login

type login struct{}

// GetLoginRecordList 获取用户登录列表
func (l *login) GetLoginRecordList(username string, page, limit int) (data *dao.LoginRecordList, err error) {
	data, err = dao.Login.GetLoginRecordList(username, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// AddLoginRecord 新增登录记录
func (l *login) AddLoginRecord(tx *gorm.DB, status int, username, loginMethod string, failedReason error, c *gin.Context) (err error) {
	// 获取登录客户端信息
	userAgent := c.Request.UserAgent()
	clientIP := c.ClientIP()

	// 数据封闭，Status=1表示成功
	loginRecord := &model.LogLogin{
		Username:   username,
		SourceIP:   clientIP,
		UserAgent:  userAgent,
		Status:     status,
		AuthMethod: loginMethod,
	}

	// 记录登录失败原因
	if status != 1 {
		loginRecord.FailedReason = failedReason.Error()
	}

	// 记录登录客户端信息
	if err := dao.Login.AddLoginRecord(tx, loginRecord); err != nil {
		return err
	}
	return nil
}
