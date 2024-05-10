package dao

import (
	"errors"
	"github.com/wonderivan/logger"
	"ops-api/global"
	"ops-api/model"
)

var Menu menu

type menu struct{}

// MenuItem 菜单项
type MenuItem struct {
	Path      string            `json:"path"`
	Component string            `json:"component"`
	Name      string            `json:"name"`
	Meta      map[string]string `json:"meta"`
	Children  []*MenuItem       `json:"children,omitempty"` // 当Children为Null时不返回，否则前端无法正确加载路由
}

// GetUserMenu 获取用户菜单
func (m *menu) GetUserMenu() (data []*MenuItem, err error) {

	var (
		menus     []*model.Menu
		menuItems []*MenuItem
	)

	// 获取一级菜单
	if err := global.MySQLClient.Find(&menus).Error; err != nil {
		logger.Error("ERROR：", err.Error())
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
			logger.Error("ERROR：", err.Error())
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
