package model

// AuthGroupPermissions 权限表
type AuthGroupPermissions struct {
	ID               int `gorm:"primaryKey;"`
	AuthGroupID      int
	AuthPermissionID int
}

func (*AuthGroupPermissions) TableName() (name string) {
	return "auth_group_permissions"
}
