package model

import (
	"gorm.io/gorm"
	"ops-api/utils"
)

// Account 账号信息表
type Account struct {
	gorm.Model
	Name         string      `json:"name" binding:"required"`
	LoginAddress string      `json:"login_address" gorm:"default:null"`
	LoginMethod  string      `json:"login_method" gorm:"default:null"`
	Username     string      `json:"username" gorm:"default:null"`
	Password     string      `json:"password" binding:"required"`
	Note         string      `json:"note" gorm:"default:null"`
	OwnerUser    *AuthUser   `gorm:"foreignKey:OwnerUserID"`
	OwnerUserID  uint        `json:"owner_user_id"`
	Users        []*AuthUser `json:"users" gorm:"many2many:account_users"`
}

func (*Account) TableName() (name string) {
	return "account"
}

// BeforeCreate 创建时对密码字段加密，仅创建时候调用
func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	cipherText, err := utils.Encrypt(a.Password)
	if err != nil {
		return err
	}
	a.Password = cipherText
	return nil
}
