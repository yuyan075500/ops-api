package model

import (
	"gorm.io/gorm"
	"time"
)

// AuthUser 用户信息表
type AuthUser struct {
	gorm.Model
	Name        string      `json:"name"`
	UserName    string      `json:"userName" gorm:"unique"`
	Avatar      *string     `json:"avatar"`
	Password    string      `json:"password"`
	PhoneNumber string      `json:"phone_number"`
	IsActive    int         `json:"is_active" gorm:"default:1"` // 1激活 0禁用
	Email       string      `json:"email"`
	LastLoginAt *time.Time  `json:"last_login_at"` // 允许为空
	MFACode     *string     `json:"mfa_code"`      // 允许为空
	UserFrom    string      `json:"user_from" gorm:"default:本地"`
	Groups      []AuthGroup `json:"groups" gorm:"many2many:auth_user_groups"`
}

func (*AuthUser) TableName() (name string) {
	return "auth_user"
}
