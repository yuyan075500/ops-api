package model

// AuthUserGroups 自定义用户与组（角色）表
type AuthUserGroups struct {
	AuthUserID  uint `gorm:"primaryKey;index:users"`
	AuthGroupID uint `gorm:"primaryKey;index:users"`
}

func (*AuthUserGroups) TableName() (name string) {
	return "auth_user_groups"
}
