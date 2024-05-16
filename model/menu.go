package model

// Menu 一级菜单
type Menu struct {
	Id        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string `json:"title"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Sort      uint   `json:"sort"`
	SubMenus  []SubMenu
}

func (*Menu) TableName() (name string) {
	return "system_menu"
}

// SubMenu 二级菜单
type SubMenu struct {
	Id        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string `json:"title"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Sort      uint   `json:"sort"`
	MenuID    uint   `json:"menu_id"`
}

func (*SubMenu) TableName() (name string) {
	return "system_sub_menu"
}
