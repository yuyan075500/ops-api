package model

import (
	"gorm.io/gorm"
	"ops-api/utils"
	"time"
)

// AuthUser 用户信息表
type AuthUser struct {
	gorm.Model
	Name        string      `json:"name"`
	Username    string      `json:"username" gorm:"unique"`
	Avatar      *string     `json:"avatar"`
	Password    string      `json:"password"`
	PhoneNumber string      `json:"phone_number"`
	IsActive    *bool       `json:"is_active" gorm:"default:true"`
	Email       string      `json:"email"`
	LastLoginAt *time.Time  `json:"last_login_at"`
	MFACode     *string     `json:"mfa_code"`
	UserFrom    string      `json:"user_from" gorm:"default:本地"`
	Groups      []AuthGroup `json:"groups" gorm:"many2many:auth_user_groups"`
}

func (*AuthUser) TableName() (name string) {
	return "auth_user"
}

// BeforeSave 新用户创建前对密码字段加密
func (u *AuthUser) BeforeSave(tx *gorm.DB) (err error) {
	cipherText, err := utils.Encrypt(u.Password)
	if err != nil {
		return err
	}
	u.Password = cipherText
	return nil
}

// CheckPassword 检查用户密码是否正确
func (u *AuthUser) CheckPassword(password string) bool {

	// 对数据库中的密码解密
	str, err := utils.Decrypt(u.Password)
	if err != nil {
		return false
	}

	// 判断密码是否相等
	if str != password {
		return false
	}

	return true
}
