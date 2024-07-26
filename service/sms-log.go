package service

import (
	"encoding/json"
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
	Result      []Result `json:"result"`
	Code        string   `json:"code"`
	Description string   `json:"description"`
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

// UserInfo 户重置密码时用户信息绑定的结构体
type UserInfo struct {
	Username    string `json:"username" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// SMSSend 发送短信
func (l *log) SMSSend(data *UserInfo, expirationTime string) (code string, err error) {

	var (
		response Response
		num      = utils.GenerateRandomNumber()
	)

	// 在数据库中查询用户是否存在
	tx := global.MySQLClient.First(&model.AuthUser{}, "username = ? AND phone_number = ?", data.Username, data.PhoneNumber)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return "", errors.New("用户名或手机号错误")
	}

	// 发送短信验证码
	resp, err := sms.Send(
		config.Conf.SMS.VerificationCode.Sender,
		config.Conf.SMS.VerificationCode.TemplateId,
		config.Conf.SMS.CallbackUrl,
		config.Conf.SMS.VerificationCode.Signature,
		[]string{data.PhoneNumber},
		[]string{data.Username, strconv.Itoa(num), expirationTime},
	)
	if err != nil {
		return "", err
	}

	// 将发送短信API请求返回的数据转换为结构体
	if err := json.Unmarshal([]byte(resp), &response); err != nil {
		return "", err
	}
	// 将短信发送的记录存入数据库
	for _, result := range response.Result {
		if response.Code != "E000510" && response.Code != "000000" {
			return "", errors.New("短信发送失败，错误码：" + response.Code)
		}

		if response.Code == "E000510" && result.Status != "000000" {
			return "", errors.New("短信发送失败，错误码：" + result.Status)
		}
		s := &model.LogSMS{
			Note:       "密码重置",
			Signature:  config.Conf.SMS.VerificationCode.Signature,
			TemplateId: config.Conf.SMS.VerificationCode.TemplateId,
			Receiver:   result.OriginTo,
			Status:     "API请求成功",
			SmsMsgId:   result.SmsMsgId,
		}
		if err := dao.Log.AddSMSRecord(s); err != nil {
			return "", err
		}
	}

	return strconv.Itoa(num), nil
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
