package dao

import (
	"errors"
	"gorm.io/gorm"
	"ops-api/global"
	"ops-api/model"
)

var Menu menu

type menu struct{}

// MenuList 返回给前端菜单列表结构体
type MenuList struct {
	Items []*model.Menu `json:"items"`
	Total int64         `json:"total"`
}

// MenuItem 菜单项
type MenuItem struct {
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Component string            `json:"component"`
	Meta      map[string]string `json:"meta"`
	Children  []*MenuItem       `json:"children,omitempty"` // 当Children为Null时不返回，否则前端无法正确加载路由
}

// GetMenuListAll 获取所有菜单
func (m *menu) GetMenuListAll() (data *MenuList, err error) {

	// 定义返回的内容
	var (
		menus []*model.Menu
		total int64
	)

	// 获取所有菜单
	tx := global.MySQLClient.Model(&model.Menu{}).
		Preload("SubMenus"). // 加载二级菜单
		Count(&total).
		Find(&menus)
	if tx.Error != nil {
		return nil, errors.New(tx.Error.Error())
	}

	return &MenuList{
		Items: menus,
		Total: total,
	}, nil
}

// GetMenuList 获取菜单列表
func (m *menu) GetMenuList(title string, page, limit int) (data *MenuList, err error) {

	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		menus []*model.Menu
		total int64
	)

	// 获取菜单列表
	tx := global.MySQLClient.Model(&model.Menu{}).
		Preload("SubMenus", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort")
		}).                                   // 加载二级菜单，指定使用sort字段进行排序
		Where("title like ?", "%"+title+"%"). // 实现过滤
		Count(&total).                        // 获取一级菜单总数
		Limit(limit).
		Offset(startSet).
		Order("sort"). // 使用sort字段进行排序
		Find(&menus)
	if tx.Error != nil {
		return nil, errors.New(tx.Error.Error())
	}

	return &MenuList{
		Items: menus,
		Total: total,
	}, nil
}

// GetUserMenu 获取用户菜单
func (m *menu) GetUserMenu() (data []*MenuItem, err error) {

	var (
		menus     []*model.Menu
		menuItems []*MenuItem
	)

	// 获取一级菜单
	if err := global.MySQLClient.Find(&menus).Error; err != nil {
		return nil, errors.New(err.Error())
	}

	for _, menu := range menus {
		// 将一级菜单模型转换为返回给前端的格式
		menuItem := &MenuItem{
			Path:      menu.Path,
			Component: menu.Component,
			Name:      menu.Name,
			Meta: map[string]string{
				"title": menu.Title,
				"icon":  menu.Icon,
			},
			Children: nil,
		}

		// 获取一级菜单对应的二级菜单
		var subMenus []*model.SubMenu
		if err := global.MySQLClient.Where("menu_id = ?", menu.Id).Find(&subMenus).Error; err != nil {
			return nil, errors.New(err.Error())
		}

		for _, subMenu := range subMenus {
			// 将二级菜单转换为返回给前端的格式
			subMenuItem := &MenuItem{
				Path:      subMenu.Path,
				Component: subMenu.Component,
				Name:      subMenu.Name,
				Meta: map[string]string{
					"title": subMenu.Title,
					"icon":  subMenu.Icon,
				},
			}

			// 将二级菜单添加到一级菜单的子菜单中
			menuItem.Children = append(menuItem.Children, subMenuItem)
		}

		// 将一级菜单添加到返回给前端的菜单列表中
		menuItems = append(menuItems, menuItem)
	}

	return menuItems, nil
}
