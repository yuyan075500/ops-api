package utils

import "math/rand"

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

// GenerateRandomNumber 生成6位非零开头的随机数字
func GenerateRandomNumber() int {
	firstDigit := rand.Intn(9) + 1
	otherDigits := rand.Intn(100000)
	result := firstDigit*100000 + otherDigits
	return result
}
