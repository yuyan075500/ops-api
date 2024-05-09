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
}

// GetUserMenu 获取用户菜单
func (m *menu) GetUserMenu() (data []*MenuItem, err error) {

	var (
		menus     []*model.Menu
		menuItems []*MenuItem
	)

	// 获取菜单列表
	if err := global.MySQLClient.Find(&menus).Error; err != nil {
		logger.Error("ERROR：", err.Error())
		return nil, errors.New(err.Error())
	}

	for _, menu := range menus {
		// 将菜单模型转换为返回给前端的格式
		menuItem := &MenuItem{
			Path:      menu.Path,
			Component: menu.Component,
			Name:      menu.Name,
			Meta: map[string]string{
				"title": menu.Title,
				"icon":  menu.Icon,
			},
		}
		menuItems = append(menuItems, menuItem)
	}

	return menuItems, nil
}
