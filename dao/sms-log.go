package dao

import (
	"errors"
	"ops-api/global"
	"ops-api/model"
)

var Log log

type log struct{}

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

// GetSMSRecordList 获取短信发送列表
func (l *log) GetSMSRecordList(receiver string, page, limit int) (data *SMSRecordList, err error) {

	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		record []*model.LogSMS
		total  int64
	)

	// 获取菜单列表
	tx := global.MySQLClient.Model(&model.LogSMS{}).
		Where("receiver like ?", "%"+receiver+"%"). // 实现过滤
		Count(&total).                              // 获取一级菜单总数
		Limit(limit).
		Offset(startSet).
		Order("id desc"). // 使用sort字段进行排序
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
func (l *log) AddSMSRecord(data *model.LogSMS) (err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// SMSCallback 短信回执
func (l *log) SMSCallback(data *Callback) (err error) {
	if err := global.MySQLClient.Model(&model.LogSMS{}).Where("sms_msg_id = ?", data.SmsMsgId).Updates(data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}
