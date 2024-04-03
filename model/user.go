package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

// AuthUser 用户表
type AuthUser struct {
	gorm.Model
	Name        string
	UserName    string `gorm:"unique"`
	Password    string
	PhoneNumber string
	IsActive    int
	Email       string
	LastLoginAt *time.Time //允许为空
	MFACode     *string    //允许为空
	UserFrom    string
	Group       []*AuthGroup `gorm:"many2many:auth_user_groups"`
}

// AuthGroup 用户分组表
type AuthGroup struct {
	gorm.Model
	Name        string            `gorm:"unique"`
	Users       []*AuthUser       `gorm:"many2many:auth_user_groups"`
	Permissions []*AuthPermission `gorm:"many2many:auth_group_permissions"`
}

// AuthPermission 权限表
type AuthPermission struct {
	gorm.Model
	Name     string       //权限中文名称
	CodeName string       //权限代号，英文名称
	Groups   []*AuthGroup `gorm:"many2many:auth_group_permissions"`
}
