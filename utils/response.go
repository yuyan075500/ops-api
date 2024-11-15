package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"reflect"
)

// excludedFields 不记录日志的字段
var excludedFields = map[string]struct{}{
	"password":      {},
	"re_password":   {},
	"client_id":     {},
	"client_secret": {},
	"certificate":   {},
	"mfa_code":      {},
	"DeletedAt":     {},
}

// FilterFields 递归过滤敏感字段
func FilterFields(data map[string]interface{}) {
	for key, value := range data {
		// 检查是否为敏感字段
		if _, ok := excludedFields[key]; ok {
			// 删除敏感字段
			delete(data, key)
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
func StructToMap(data interface{}) (interface{}, error) {
	// 如果是结构体，直接转换为map
	if reflect.TypeOf(data).Kind() != reflect.Slice {
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

	// 如果是切片类型，处理切片中的每一项
	var result []map[string]interface{}
	v := reflect.ValueOf(data)
	// 遍历切片中的每一项
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		itemData, err := StructToMap(item)
		if err != nil {
			return nil, err
		}
		result = append(result, itemData.(map[string]interface{}))
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
// SendCreateOrUpdateResponse 创建或更新请求的响应
func SendCreateOrUpdateResponse(c *gin.Context, code int, msg string, data interface{}) {

	var responseData interface{}
	if data != nil {
		// 将结构体转换为 map 或 []map
		responseData, _ = StructToMap(data)

		// 检查 responseData 的类型
		switch v := responseData.(type) {
		case map[string]interface{}:
			// 如果是单个对象，直接过滤敏感信息
			FilterFields(v)
		case []map[string]interface{}:
			// 如果是数组，遍历每个元素并过滤敏感信息
			for _, item := range v {
				FilterFields(item)
			}
		}
	}

	response := map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": responseData,
	}
	c.Set("response", response)
	c.JSON(200, response)
}
