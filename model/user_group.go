package model

// AuthUserGroups 自定义用户与组（角色）表
type AuthUserGroups struct {
	ID          uint `gorm:"primaryKey;"`
	AuthUserID  uint
	AuthGroupID uint
}

func (*AuthUserGroups) TableName() (name string) {
	return "auth_user_groups"
}
