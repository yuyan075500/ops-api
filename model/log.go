package model

import (
	"gorm.io/gorm"
)

// LogSMS 短信发送日志表
type LogSMS struct {
	gorm.Model
	Note       string  `json:"note"`
	Signature  string  `json:"signature"`
	TemplateId string  `json:"template_id"`
	Receiver   string  `json:"receiver"`
	Status     string  `json:"status"`
	ErrorCode  *string `json:"error_code"`
	SmsMsgId   string  `json:"sms_msg_id"`
}

// LogUserLogin 用户登录日志表
type LogUserLogin struct {
	gorm.Model
	Username   string `json:"username"`
	SourceIP   string `json:"source_ip"`
	UserAgent  string `json:"user_agent"`
	Status     string `json:"status"`
	AuthMethod string `json:"auth_method"`
	City       string `json:"city"`
}
