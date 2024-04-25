package check

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// PhoneNumberCheck 手机号验证，用于结构体中手机号检验
func PhoneNumberCheck(n validator.FieldLevel) bool {

	// 获取手机号
	phoneNumber := n.Field().String()

	// 手机号码验证规则
	ruler := "^1[345789]{1}\\d{9}$"
	reg := regexp.MustCompile(ruler)

	return reg.MatchString(phoneNumber)
}
