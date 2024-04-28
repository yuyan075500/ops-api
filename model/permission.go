package model

// AuthPermission 基于CasBin权限表
type AuthPermission struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	PType  string `json:"ptype"` // 权限类型
	RoleId uint   `json:"v0"`    // 角色ID
	Path   string `json:"v1"`    // 路径
	Method string `json:"v2"`    // 方法
}

func (*AuthPermission) TableName() (name string) {
	return "auth_permission"
}

// AddPolicy 添加权限
func (auth *AuthPermission) AddPolicy() error {
	return nil
}
