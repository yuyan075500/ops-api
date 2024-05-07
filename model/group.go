package model

import "gorm.io/gorm"

// AuthGroup 用户组、角色信息表
type AuthGroup struct {
	gorm.Model
	Name        string      `json:"name" gorm:"unique"`
	IsRoleGroup bool        `json:"is_role_group" gorm:"default:false"`
	Users       []*AuthUser `json:"users" gorm:"many2many:auth_user_groups"`
}

func (*AuthGroup) TableName() (name string) {
	return "auth_group"
}
