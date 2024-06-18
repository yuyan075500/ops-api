package dao

import (
	"errors"
	"ops-api/global"
	"ops-api/model"
)

var Log log

type log struct{}

type Callback struct {
	Status    string `json:"status"`
	SmsMsgId  string `json:"smsMsgId"`
	ErrorCode string `json:"code"`
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
