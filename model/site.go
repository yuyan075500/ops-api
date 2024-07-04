package model

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
	Icon         *string     `json:"icon"`
	Address      string      `json:"address"`
	AllOpen      bool        `json:"all_open"`
	Description  string      `json:"description"`
	SSO          bool        `json:"sso"`
	SSOType      string      `json:"sso_type"`
	ClientId     string      `json:"client_id"`
	ClientSecret string      `json:"client_secret"`
	CallbackUrl  string      `json:"callback_url"`
	SiteGroupID  uint        `json:"site_group_id"`
	Users        []*AuthUser `json:"users" gorm:"many2many:auth_user_sites"`
}

func (*Site) TableName() (name string) {
	return "site"
}
