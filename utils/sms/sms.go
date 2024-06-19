package sms

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"ops-api/config"
	"ops-api/utils/sms/core"
)

// Send 发送短信（华为云）
//
// Param:
//
//	sender: 国内短信签名通道号
//	templateId: 短信模板ID
//	statusCallBack: 短信发送状态回调地址，为空则表示不接收
//	signature: 短信签名名称
//	receiver: 接口短信的电话号码，多个号码之间使用英文逗号分隔
//	templateParas: 模板参数
func Send(sender, templateId, statusCallBack, signature string, receiver, templateParas []string) (resp string, err error) {

	// 短信初始化，支持一次发送多条，请参考官方文档
	sms := initDiffSms(receiver, templateId, templateParas, signature)

	// 创建签名对象
	appInfo := core.Signer{
		Key:    config.Conf.SMS.AppKey,
		Secret: config.Conf.SMS.AppSecret,
	}

	// 请求接口地址
	apiAddress := config.Conf.SMS.URL

	// 构造请求体
	body := buildRequestBody(sender, []map[string]interface{}{sms}, statusCallBack)

	// 发送短信请求
	resp, err = post(apiAddress, []byte(body), appInfo)
	if err != nil {
		return "", err
	}

	// 返回API请求响应
	return resp, nil
}

// buildRequestBody 构造发送短信的请求体
func buildRequestBody(sender string, item []map[string]interface{}, statusCallBack string) []byte {
	body := make(map[string]interface{})
	body["smsContent"] = item
	body["from"] = sender
	if statusCallBack != "" {
		body["statusCallback"] = statusCallBack
	}
	res, _ := json.Marshal(body)
	return res
}

func initDiffSms(receiver []string, templateId string, templateParas []string, signature string) map[string]interface{} {
	diffSms := make(map[string]interface{})
	diffSms["to"] = receiver
	diffSms["templateId"] = templateId
	if templateParas != nil && len(templateParas) > 0 {
		diffSms["templateParas"] = templateParas
	}
	if signature != "" {
		diffSms["signature"] = signature
	}
	return diffSms
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
	req.Header.Add("Content-Type", "application/json")

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
