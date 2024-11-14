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
	Redirect  string            `json:"redirect,omitempty"`
	Children  []*MenuItem       `json:"children,omitempty"` // 当Children为Null时不返回，否则前端无法正确加载路由
}

// GetMenuListAll 获取所有菜单（权限分配）
func (m *menu) GetMenuListAll() (data *MenuList, err error) {

	// 定义返回的内容
	var (
		menus []*model.Menu
		total int64
	)

	// 获取所有菜单
	tx := global.MySQLClient.Model(&model.Menu{}).
		Preload("SubMenus", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort")
		}). // 加载二级菜单
		Count(&total).
		Order("sort").
		Find(&menus)
	if tx.Error != nil {
		return nil, err
	}

	return &MenuList{
		Items: menus,
		Total: total,
	}, nil
}

// GetMenuList 获取菜单列表（表格中展示）
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
		return nil, err
	}

	return &MenuList{
		Items: menus,
		Total: total,
	}, nil
}

// GetUserMenu 获取用户有菜单（用户登录）
func (m *menu) GetUserMenu(tx *gorm.DB, username string) (data []*MenuItem, err error) {

	var (
		menus     []*model.Menu
		menuItems []*MenuItem
	)

	// 获取一级菜单
	if err := tx.Order("sort").Find(&menus).Error; err != nil {
		return nil, err
	}

	for _, menu := range menus {
		// 判断用户是否拥有该菜单权限
		ok, _ := global.CasBinServer.Enforce(username, menu.Name, "read")
		if ok {
			// 将一级菜单模型转换为返回给前端的格式
			menuItem := &MenuItem{
				Path:      menu.Path,
				Component: menu.Component,
				Name:      menu.Name,
				Redirect:  menu.Redirect,
				Meta: map[string]string{
					"title": menu.Title,
					"icon":  menu.Icon,
				},
				Children: nil,
			}

			// 获取一级菜单对应的二级菜单
			var subMenus []*model.SubMenu
			if err := tx.Where("menu_id = ?", menu.Id).Order("sort").Find(&subMenus).Error; err != nil {
				return nil, err
			}
			for _, subMenu := range subMenus {
				// 判断用户是否拥有该菜单权限
				ok, _ := global.CasBinServer.Enforce(username, subMenu.Name, "read")
				if ok {
					// 将二级菜单转换为返回给前端的格式
					subMenuItem := &MenuItem{
						Path:      subMenu.Path,
						Component: subMenu.Component,
						Name:      subMenu.Name,
						Redirect:  subMenu.Redirect,
						Meta: map[string]string{
							"title": subMenu.Title,
							"icon":  subMenu.Icon,
						},
					}

					// 将二级菜单添加到一级菜单的子菜单中
					menuItem.Children = append(menuItem.Children, subMenuItem)
				}
			}

			// 将一级菜单添加到返回给前端的菜单列表中
			menuItems = append(menuItems, menuItem)
		}
	}

	return menuItems, nil
}

// GetMenuTitle 根据菜单Name获取Title
func (m *menu) GetMenuTitle(menuName string) (title *string, err error) {
	var (
		menu    model.Menu
		subMenu model.SubMenu
	)

	// 在一级菜单中根据Name获取Title
	tx := global.MySQLClient.Where("name = ?", menuName).First(&menu)

	// 如果一级菜单没有找到对应的记录则在二级菜单中继续查找
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		tx := global.MySQLClient.Where("name = ?", menuName).First(&subMenu)
		if tx.Error != nil {
			return nil, err
		}
		return &subMenu.Title, nil
	}

	if tx.Error != nil {
		return nil, err
	}

	return &menu.Title, nil
}
