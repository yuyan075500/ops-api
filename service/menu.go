package service

import "ops-api/dao"

var Menu menu

type menu struct{}

// GetMenuList 获取菜单列表
func (m *menu) GetMenuList(page, limit int) (data *dao.MenuList, err error) {
	data, err = dao.Menu.GetMenuList(page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetUserMenu 获取用户菜单
func (m *menu) GetUserMenu() (data []*dao.MenuItem, err error) {
	data, err = dao.Menu.GetUserMenu()
	if err != nil {
		return nil, err
	}

	return data, nil
}
