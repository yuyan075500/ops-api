package service

import (
	"errors"
	"gorm.io/gorm"
	"net/url"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils"
	"ops-api/utils/sms"
	"strconv"
	"strings"
)

var Log log

type log struct{}

// Response 短信API发送请求返回的数据
type Response struct {
	Result      []Result `json:"result"`      // 华为云
	Code        string   `json:"code"`        // 华为云
	Description string   `json:"description"` // 华为云
	Body        Body     `json:"body"`        // 阿里云
	StatusCode  int      `json:"statusCode"`  // 阿里云
}
type Body struct {
	BizId     string `json:"BizId"`
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	RequestId string `json:"RequestId"`
}

type Result struct {
	Total      int    `json:"total"`
	OriginTo   string `json:"originTo"`
	CreateTime string `json:"createTime"`
	From       string `json:"from"`
	SmsMsgId   string `json:"smsMsgId"`
	CountryId  string `json:"countryId"`
	Status     string `json:"status"`
}

// SMSSend 发送短信
func (l *log) SMSSend(data *sms.ResetPassword) (string, error) {

	// 定义验证码
	var code = strconv.Itoa(utils.GenerateRandomNumber())

	// 查询用户是否存在
	tx := global.MySQLClient.First(&model.AuthUser{}, "username = ? AND phone_number = ?", data.Username, data.PhoneNumber)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return "", errors.New("用户名或手机号错误")
	}

	// 获取短信服务商
	smsSender := sms.GetSMSSender()
	if smsSender == nil {
		return "", errors.New("不支持的短信服务提供商")
	}

	// 发送短信
	resp, err := smsSender.SendSMS(data, code)
	if err != nil {
		return "", err
	}

	// 处理响应并获取短信唯一标识
	smsMsgId, err := smsSender.ProcessResponse(resp)
	if err != nil {
		return "", err
	}

	// 记录短信发送日志
	smsLog := &model.LogSMS{
		Note:       "密码重置",
		Signature:  config.Conf.SMS.ResetPassword.Signature,
		TemplateId: config.Conf.SMS.ResetPassword.TemplateId,
		Receiver:   data.PhoneNumber,
		Status:     "API请求成功",
		SmsMsgId:   smsMsgId,
	}

	if err := dao.Log.AddSMSRecord(smsLog); err != nil {
		return "", err
	}

	return code, nil
}

// GetSMSRecordList 获取发送短信列表
func (l *log) GetSMSRecordList(receiver string, page, limit int) (data *dao.SMSRecordList, err error) {
	data, err = dao.Log.GetSMSRecordList(receiver, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SMSCallback 短信回调
func (l *log) SMSCallback(data string) error {

	// 对回调返回的数据进行处理
	ss, _ := url.QueryUnescape(data)
	params := strings.Split(ss, "&")
	keyValues := make(map[string]string)
	for i := range params {
		temp := strings.Split(params[i], "=")
		keyValues[temp[0]] = temp[1]
	}

	// 将数据与结构体进行绑定
	callback := &dao.Callback{
		Status:    "接收成功",
		SmsMsgId:  keyValues["smsMsgId"],
		ErrorCode: "",
	}
	if keyValues["status"] != "DELIVRD" {
		callback.Status = "发送失败"
		callback.ErrorCode = keyValues["status"]
	}

	// 将回调数据写入数据库
	if err := dao.Log.SMSCallback(callback); err != nil {
		return err
	}

	return nil
}
