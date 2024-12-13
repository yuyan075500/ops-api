package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"ops-api/utils"
	"reflect"
)

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

// Response 普通请求响应
func Response(c *gin.Context, code int, msg string) {

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

// CreateOrUpdateResponse 创建或更新请求的响应
func CreateOrUpdateResponse(c *gin.Context, code int, msg string, data interface{}) {

	var responseData interface{}
	if data != nil {
		// 将结构体转换为 map 或 []map
		responseData, _ = StructToMap(data)

		// 检查 responseData 的类型
		switch v := responseData.(type) {
		case map[string]interface{}:
			// 如果是单个对象，直接过滤敏感信息
			utils.FilterFields(v)
		case []map[string]interface{}:
			// 如果是数组，遍历每个元素并过滤敏感信息
			for _, item := range v {
				utils.FilterFields(item)
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
