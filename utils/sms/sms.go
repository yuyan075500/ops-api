package sms

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/model"
	"ops-api/utils/sms/core"
)

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

// Send 发送短信（华为云）
//
// Param:
//
//		sender: 国内短信签名通道号
//		receiver: 接口短信的电话号码，多个号码之间使用英文逗号分隔
//		templateId: 短信模板ID
//		templateParas: 模板参数
//		statusCallBack: 短信发送状态回调地址，为空则表示不接收
//		signature: 短信签名名称
//	 note: 标识短信用途
func Send(sender, receiver, templateId, templateParas, statusCallBack, signature, note string) error {

	// 创建签名对象
	appInfo := core.Signer{
		Key:    config.Conf.SMS.AppKey,
		Secret: config.Conf.SMS.AppSecret,
	}

	// 请求接口地址
	apiAddress := config.Conf.SMS.URL

	// 构造请求体
	body := buildRequestBody(sender, receiver, templateId, templateParas, statusCallBack, signature)

	// 发送短信请求
	resp, err := post(apiAddress, []byte(body), appInfo)
	if err != nil {
		return err
	}

	// 解析API请求返回的数据
	var response Response
	if err := json.Unmarshal([]byte(resp), &response); err != nil {
		return err
	}

	// 在数据库中创建短信发送记录
	for _, result := range response.Result {
		if response.Code != "E000510" && response.Code != "000000" {
			return errors.New("短信发送失败，错误码：" + response.Code)
		}

		if response.Code == "E000510" && result.Status != "000000" {
			return errors.New("短信发送失败，错误码：" + result.Status)
		}
		sms := &model.LogSMS{
			Note:       note,
			Signature:  signature,
			TemplateId: templateId,
			Receiver:   receiver,
			Status:     "API请求成功",
			SmsMsgId:   result.SmsMsgId,
		}
		if err := dao.Log.AddSMSRecord(sms); err != nil {
			return err
		}
	}

	return nil
}

// buildRequestBody 构造发送短信的请求体
func buildRequestBody(sender, receiver, templateId, templateParas, statusCallBack, signature string) string {
	param := "from=" + url.QueryEscape(sender) + "&to=" + url.QueryEscape(receiver) + "&templateId=" + url.QueryEscape(templateId)
	if templateParas != "" {
		param += "&templateParas=" + url.QueryEscape(templateParas)
	}
	if statusCallBack != "" {
		param += "&statusCallback=" + url.QueryEscape(statusCallBack)
	}
	if signature != "" {
		param += "&signature=" + url.QueryEscape(signature)
	}
	return param
}

// post 发送短信请求
func post(url string, param []byte, appInfo core.Signer) (string, error) {
	if param == nil || appInfo == (core.Signer{}) {
		return "", nil
	}

	// 代码样例为了简便，设置了不进行证书校验，请在商用环境自行开启证书校验。
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(param))
	if err != nil {
		return "", err
	}

	// 对请求增加内容格式，固定头域
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 对请求进行HMAC算法签名，并将签名结果设置到Authorization头域。
	if err := appInfo.Sign(req); err != nil {
		return "", err
	}

	// 发送短信请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// 获取短信响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
