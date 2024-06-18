package service

import (
	"fmt"
	"net/url"
	"ops-api/dao"
	"strings"
)

var Log log

type log struct{}

// GetSMSRecordList 获取菜单列表
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
	fmt.Println(callback)
	if err := dao.Log.SMSCallback(callback); err != nil {
		return err
	}

	return nil
}
