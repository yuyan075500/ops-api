package model

import "gorm.io/gorm"

// AuthPermission 权限表
type AuthPermission struct {
	gorm.Model
	Name     string `json:"name"`      // 权限中文名称
	CodeName string `json:"code_name"` // 权限代号，英文名称
}

func (*AuthPermission) TableName() (name string) {
	return "auth_permission"
}
