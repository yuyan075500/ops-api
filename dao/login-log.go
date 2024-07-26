package dao

import (
	"errors"
	"gorm.io/gorm"
	"ops-api/global"
	"ops-api/model"
)

var Login login

type login struct{}

// LoginRecordList 返回给前端列表结构体
type LoginRecordList struct {
	Items []*model.LogLogin `json:"items"`
	Total int64             `json:"total"`
}

// LoginRecord 登录记录结构体
type LoginRecord struct {
	Username   string `json:"username"`
	SourceIP   string `json:"source_ip"`
	UserAgent  string `json:"user_agent"`
	Status     string `json:"status"`
	AuthMethod string `json:"auth_method"`
}

// GetLoginRecordList 获取用户登录列表
func (l *login) GetLoginRecordList(username string, page, limit int) (data *LoginRecordList, err error) {

	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		record []*model.LogLogin
		total  int64
	)

	// 获取菜单列表
	tx := global.MySQLClient.Model(&model.LogLogin{}).
		Where("username like ?", "%"+username+"%").
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&record)
	if tx.Error != nil {
		return nil, errors.New(tx.Error.Error())
	}

	return &LoginRecordList{
		Items: record,
		Total: total,
	}, nil
}

// AddLoginRecord 新增用户登录记录
func (l *login) AddLoginRecord(tx *gorm.DB, data *model.LogLogin) (err error) {
	if err := tx.Create(&data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}
