package utils

// Contains 查询字符串在一个列表中是否存在
func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// MapToJson Map转json
//func MapToJson(data interface{}) string {
//	byteStr, _ := json.Marshal(data)
//	return string(byteStr)
//}
