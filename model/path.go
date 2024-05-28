package model

// SystemPath 系统API接口
type SystemPath struct {
	Id          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"unique"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	MenuName    string `json:"menu_name"`
	Description string `json:"description"`
}

func (*SystemPath) TableName() (name string) {
	return "system_path"
}
