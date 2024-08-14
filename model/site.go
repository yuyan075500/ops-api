package model

import (
	"gorm.io/gorm"
	"ops-api/utils"
)

// SiteGroup 站点分组
type SiteGroup struct {
	ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name"`
	Sites []Site
}

func (*SiteGroup) TableName() (name string) {
	return "site_group"
}

// Site 站点
type Site struct {
	ID           uint        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string      `json:"name"`
	Icon         *string     `json:"icon" gorm:"default:null"`
	Address      string      `json:"address"`
	AllOpen      bool        `json:"all_open" gorm:"default:false"`
	Description  string      `json:"description"`
	SSO          bool        `json:"sso"`
	SSOType      uint        `json:"sso_type" gorm:"default:null"`
	ClientId     string      `json:"client_id"`                        // OAuth2.0 ClientID
	ClientSecret string      `json:"client_secret"`                    // OAuth2.0 ClientSecret
	CallbackUrl  string      `json:"callback_url" gorm:"default:null"` // OAuth2.0 And CAS3.0 Client CallbackUrl
	EntityId     string      `json:"entity_id" gorm:"default:null"`    // SAML2.0 SP EntityID
	Certificate  string      `json:"certificate" gorm:"default:null"`  // SAML2.0 SP Certificate
	SiteGroupID  uint        `json:"site_group_id"`
	Users        []*AuthUser `json:"users" gorm:"many2many:site_users"`
}

func (*Site) TableName() (name string) {
	return "site"
}

// BeforeCreate 创建新的站点时生成ClientId和ClientSecret
func (s *Site) BeforeCreate(tx *gorm.DB) (err error) {
	s.ClientId = utils.GenerateRandomString(16)
	s.ClientSecret = utils.GenerateRandomString(32)
	return nil
}
