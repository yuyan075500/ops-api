package dao

import (
	"errors"
	"gorm.io/gorm"
	"ops-api/global"
	"ops-api/model"
)

var Audit audit

type audit struct{}

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

// SMSRecordList 返回给前端短信发送列表结构体
type SMSRecordList struct {
	Items []*model.LogSMS `json:"items"`
	Total int64           `json:"total"`
}

// Callback 短信回执结构体
type Callback struct {
	Status    string `json:"status"`
	SmsMsgId  string `json:"smsMsgId"`
	ErrorCode string `json:"code"`
}

// GetLoginRecordList 获取系统登录记录
func (a *audit) GetLoginRecordList(name string, page, limit int) (data *LoginRecordList, err error) {

	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		record []*model.LogLogin
		total  int64
	)

	// 获取菜单列表
	tx := global.MySQLClient.Model(&model.LogLogin{}).
		Where("username like ? OR source_ip like ? OR user_agent like ? OR city like ?", "%"+name+"%", "%"+name+"%", "%"+name+"%", "%"+name+"%").
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

// GetSMSRecordList 获取短信发送记录
func (a *audit) GetSMSRecordList(receiver string, page, limit int) (data *SMSRecordList, err error) {

	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		record []*model.LogSMS
		total  int64
	)

	// 获取菜单列表
	tx := global.MySQLClient.Model(&model.LogSMS{}).
		Where("receiver like ?", "%"+receiver+"%").
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&record)
	if tx.Error != nil {
		return nil, errors.New(tx.Error.Error())
	}

	return &SMSRecordList{
		Items: record,
		Total: total,
	}, nil
}

// AddSMSRecord 新增短信发送记录
func (a *audit) AddSMSRecord(data *model.LogSMS) (err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// GetSendDetail 获取短信发送详情
func (a *audit) GetSendDetail(conditions interface{}) (*model.LogSMS, error) {

	var sms model.LogSMS

	if err := global.MySQLClient.Where(conditions).First(&sms).Error; err != nil {
		return nil, err
	}

	return &sms, nil
}

// SMSCallback 短信回执
func (a *audit) SMSCallback(data *Callback) (err error) {
	if err := global.MySQLClient.Model(&model.LogSMS{}).Where("sms_msg_id = ?", data.SmsMsgId).Updates(data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// AddLoginRecord 新增用户登录记录
func (a *audit) AddLoginRecord(tx *gorm.DB, data *model.LogLogin) (err error) {
	if err := tx.Create(&data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}
