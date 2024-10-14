package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/model"
	"ops-api/utils/sms"
)

var Login login

type login struct{}

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
func (l *login) GetSMSReceipt(smsId int) (err error) {

	// 华为云不需要
	if config.Conf.SMS.Provider != "aliyun" {
		return nil
	}

	// 定义匹配条件
	conditions := map[string]interface{}{
		"id": smsId,
	}

	// 查找短信记录
	smsRecord, err := dao.Log.GetSendDetail(conditions)
	if err != nil {
		return err
	}

	// 获取短信回执
	date := smsRecord.CreatedAt
	receipt, err := sms.GetSMSReceipt(smsRecord.Receiver, smsRecord.SmsMsgId, date.Format("20060102"))
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
	if err := dao.Log.SMSCallback(callback); err != nil {
		return err
	}

	return nil
}

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
	if err := dao.Login.AddLoginRecord(tx, loginRecord); err != nil {
		return err
	}
	return nil
}
