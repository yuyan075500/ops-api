package model

import "gorm.io/gorm"

// AuthGroup 用户组、角色信息表
type AuthGroup struct {
	gorm.Model
	Name        string           `json:"name" gorm:"unique"`
	Permissions []AuthPermission `json:"permissions" gorm:"many2many:auth_group_permissions"`
}

func (*AuthGroup) TableName() (name string) {
	return "auth_group"
}
