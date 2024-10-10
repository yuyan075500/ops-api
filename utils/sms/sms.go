package sms

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"ops-api/config"
	"ops-api/utils/sms/core"
	"strings"
)

// Response 短信发送结果返回结构体
type Response struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// Send 发送短信（华为云）
//
// Param:
//
//	sender: 国内短信签名通道号
//	templateId: 短信模板ID
//	statusCallBack: 短信发送状态回调地址，为空则表示不接收
//	signature: 短信签名名称
//	receiver: 接口短信的电话号码
//	templateParas: 模板参数
func Send(sender, templateId, statusCallBack, signature, receiver string, templateParas int) (resp string, err error) {

	var rsp Response

	// 创建签名对象
	appInfo := core.Signer{
		Key:    config.Conf.SMS.AppKey,
		Secret: config.Conf.SMS.AppSecret,
	}

	// 请求接口地址
	apiAddress := config.Conf.SMS.URL

	// 构造参数
	paras := fmt.Sprintf("[\"%d\"]", templateParas)

	// 构造请求体
	body := buildRequestBody(sender, receiver, templateId, paras, statusCallBack, signature)

	// 发送短信请求
	resp, err = post(apiAddress, []byte(body), appInfo)
	if err != nil {
		return "", err
	}

	// 将JSON字符串解析为结构体
	err = json.Unmarshal([]byte(resp), &rsp)
	if err != nil {
		return "", err
	}

	// 返回API请求响应
	if rsp.Code != "000000" {
		return "", errors.New(fmt.Sprintf("%s: %s", rsp.Code, rsp.Description))
	}
	return resp, nil
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
	appInfo.Sign(req)

	// 发送短信请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// 获取短信响应
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// 格式化模板参数函数
func formatTemplateParas(vars []string) string {
	// 使用转义后的双引号包裹每个变量
	for i, v := range vars {
		vars[i] = fmt.Sprintf("\\\"%s\\\"", v)
	}

	// 将所有变量拼接成一个字符串，并在最外层加上双引号
	return fmt.Sprintf("\"[%s]\"", strings.Join(vars, ","))
}
