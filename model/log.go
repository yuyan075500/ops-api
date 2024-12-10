package model

import (
	"github.com/oschwald/geoip2-golang"
	"gorm.io/gorm"
	"net"
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

func (*LogSMS) TableName() (name string) {
	return "log_sms"
}

// LogLogin 用户登录日志表
type LogLogin struct {
	gorm.Model
	Username     string `json:"username"`
	SourceIP     string `json:"source_ip"`
	UserAgent    string `json:"user_agent"`
	Status       int    `json:"status"`
	FailedReason string `json:"failed_reason"`
	AuthMethod   string `json:"auth_method"`
	City         string `json:"city"`
	Application  string `json:"application"`
}

func (*LogLogin) TableName() (name string) {
	return "log_login"
}

type LogOplog struct {
	gorm.Model
	Username      string `json:"username"`
	Endpoint      string `json:"endpoint"`
	Method        string `json:"method"`
	RequestParams string `json:"request_params"`
	ResponseData  string `json:"response_data"`
	ClientIP      string `json:"client_ip"`
	UserAgent     string `json:"user_agent"`
}

func (*LogOplog) TableName() (name string) {
	return "log_oplog"
}

// BeforeCreate 分析IP地址来源
func (l *LogLogin) BeforeCreate(tx *gorm.DB) (err error) {
	// 打开GeoLite2 City数据库
	db, err := geoip2.Open("config/GeoLite2-City.mmdb")
	if err != nil {
		return err
	}
	defer db.Close()

	// 要查询的IP地址
	ip := net.ParseIP(l.SourceIP)

	// 查询IP地址的城市信息
	record, err := db.City(ip)
	if err != nil {
		return err
	}

	city := record.City.Names["zh-CN"]
	if city == "" {
		l.City = "未知"
	} else {
		l.City = record.City.Names["zh-CN"]
	}

	return nil
}
