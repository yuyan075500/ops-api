package check

import (
	"errors"
	"fmt"
	"unicode"
)

// PasswordCheck 检查密码复杂度，用于检查密码复杂度
func PasswordCheck(password string) error {
	fmt.Println(password)

	if len(password) < 8 {
		return errors.New("密码长度不能少于8个字符")
	}

	uppercaseCount := 0
	lowercaseCount := 0
	digitCount := 0
	specialCharCount := 0

	for _, char := range password {
		if unicode.IsUpper(char) {
			uppercaseCount++
		}
		if unicode.IsLower(char) {
			lowercaseCount++
		}
		if unicode.IsDigit(char) {
			digitCount++
		}
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			specialCharCount++
		}
	}

	if uppercaseCount < 2 {
		return errors.New("至少需要包含两个大写字母")
	}
	if lowercaseCount < 2 {
		return errors.New("至少需要包含两个小写字母")
	}
	if digitCount < 2 {
		return errors.New("至少需要包含两个数字")
	}
	if specialCharCount < 2 {
		return errors.New("至少需要包含两个特殊字符")
	}

	return nil
}
