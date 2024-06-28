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

// GenerateRandomString 生成随机字符串
func GenerateRandomString(n int) string {

	// 指定随机字符串包含的字符集
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
