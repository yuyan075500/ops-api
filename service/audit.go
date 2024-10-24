package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/model"
	messages "ops-api/utils/sms"
)

var Audit audit

type audit struct{}

type Result struct {
	Total      int    `json:"total"`
	OriginTo   string `json:"originTo"`
	CreateTime string `json:"createTime"`
	From       string `json:"from"`
	SmsMsgId   string `json:"smsMsgId"`
	CountryId  string `json:"countryId"`
	Status     string `json:"status"`
}

// AliyunSMSReceipt 阿里云短信回执
type AliyunSMSReceipt struct {
	Body       ResponseBody      `json:"body"`
	Headers    map[string]string `json:"headers"`
	StatusCode int               `json:"statusCode"`
}
type SmsSendDetailDTOs struct {
	SmsSendDetailDTO []SmsSendDetailDTO `json:"SmsSendDetailDTO"`
}
type SmsSendDetailDTO struct {
	Content      string `json:"Content"`
	ErrCode      string `json:"ErrCode"`
	PhoneNum     string `json:"PhoneNum"`
	ReceiveDate  string `json:"ReceiveDate"`
	SendDate     string `json:"SendDate"`
	SendStatus   int    `json:"SendStatus"`
	TemplateCode string `json:"TemplateCode"`
}
type ResponseBody struct {
	Code              string            `json:"Code"`
	Message           string            `json:"Message"`
	RequestId         string            `json:"RequestId"`
	SmsSendDetailDTOs SmsSendDetailDTOs `json:"SmsSendDetailDTOs"`
	TotalCount        int               `json:"TotalCount"`
}

// GetSMSReceipt 获取短信回执
func (a *audit) GetSMSReceipt(smsId int) (err error) {

	// 华为云不需要
	if config.Conf.SMS.Provider != "aliyun" {
		return nil
	}

	// 定义匹配条件
	conditions := map[string]interface{}{
		"id": smsId,
	}

	// 查找短信记录
	smsRecord, err := dao.Audit.GetSendDetail(conditions)
	if err != nil {
		return err
	}

	// 获取短信回执
	date := smsRecord.CreatedAt
	receipt, err := messages.GetSMSReceipt(smsRecord.Receiver, smsRecord.SmsMsgId, date.Format("20060102"))
	if err != nil {
		return err
	}

	// 对数据进行解析
	var response AliyunSMSReceipt
	err = json.Unmarshal([]byte(*receipt), &response)
	if err != nil {
		return
	}

	// 处理回执信息
	callback := &dao.Callback{
		SmsMsgId:  smsRecord.SmsMsgId,
		ErrorCode: "",
	}
	if response.Body.Code == "OK" {
		// 获取短信回执内容
		for _, detail := range response.Body.SmsSendDetailDTOs.SmsSendDetailDTO {
			if detail.SendStatus == 1 {
				callback.Status = "等待回执"
			}

			if detail.SendStatus == 2 {
				callback.Status = "发送失败"
			}

			if detail.SendStatus == 3 {
				callback.Status = "接收成功"
			}

		}
	} else {
		callback.Status = "发送失败"
		callback.ErrorCode = response.Body.Code
	}

	// 将回调数据写入数据库
	if err := dao.Audit.SMSCallback(callback); err != nil {
		return err
	}

	return nil
}

// GetSMSRecordList 获取短信发送记录
func (a *audit) GetSMSRecordList(receiver string, page, limit int) (data *dao.SMSRecordList, err error) {
	data, err = dao.Audit.GetSMSRecordList(receiver, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetLoginRecordList 获取系统登录记录
func (a *audit) GetLoginRecordList(name string, page, limit int) (data *dao.LoginRecordList, err error) {
	data, err = dao.Audit.GetLoginRecordList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// AddLoginRecord 新增系统登录记录
func (a *audit) AddLoginRecord(tx *gorm.DB, status int, username, loginMethod string, failedReason error, c *gin.Context) (err error) {
	// 获取客户端Agent
	userAgent := c.Request.UserAgent()
	// 获取客户端IP
	clientIP := c.ClientIP()

	// 数据封装，Status=1表示成功
	loginRecord := &model.LogLogin{
		Username:   username,
		SourceIP:   clientIP,
		UserAgent:  userAgent,
		Status:     status,
		AuthMethod: loginMethod,
	}

	// 如果是登录失败，则记录登录失败原因
	if status != 1 {
		loginRecord.FailedReason = failedReason.Error()
	}

	// 记录登录客户端信息
	if err := dao.Audit.AddLoginRecord(tx, loginRecord); err != nil {
		return err
	}
	return nil
}
