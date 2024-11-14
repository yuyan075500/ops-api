package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

// excludedFields 不记录日志的字段
var excludedFields = map[string]struct{}{
	"password":      {},
	"client_id":     {},
	"client_secret": {},
	"certificate":   {},
	"mfa_code":      {},
}

// FilterFields 递归过滤敏感字段
func FilterFields(data map[string]interface{}) {
	for key, value := range data {
		// 检查是否为敏感字段
		if _, ok := excludedFields[key]; ok {
			// 将敏感字段置为nil
			data[key] = nil
		} else if nestedMap, ok := value.(map[string]interface{}); ok {
			// 递归处理嵌套 map
			FilterFields(nestedMap)
		} else if nestedArray, ok := value.([]interface{}); ok {
			// 递归处理数组中的map
			for _, item := range nestedArray {
				if itemMap, ok := item.(map[string]interface{}); ok {
					FilterFields(itemMap)
				}
			}
		}
	}
}

// StructToMap 将结构体转换为：map[string]interface{}
func StructToMap(data interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// SendResponse 普通请求响应
func SendResponse(c *gin.Context, code int, msg string) {

	// 打印错误日志
	if code != 0 {
		logger.Error("ERROR：" + msg)
	}

	// 设置响应
	response := map[string]interface{}{
		"code": code,
		"msg":  msg,
	}
	c.Set("response", response)
	c.JSON(200, response)
}

// SendCreateOrUpdateResponse 创建或更新请求的响应
func SendCreateOrUpdateResponse(c *gin.Context, code int, msg string, data interface{}) {

	// 将结构体转换为map[string]interface{}
	responseData, _ := StructToMap(data)
	// 过滤敏感信息
	FilterFields(responseData)

	response := map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": responseData,
	}
	c.Set("response", response)
	c.JSON(200, response)
}
