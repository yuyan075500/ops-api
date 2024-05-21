package service

import "ops-api/dao"

var Menu menu

type menu struct{}

// GetMenuListAll 获取菜单列表（权限分配）
func (m *menu) GetMenuListAll() (data *dao.MenuList, err error) {
	data, err = dao.Menu.GetMenuListAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetMenuList 获取菜单列表
func (m *menu) GetMenuList(title string, page, limit int) (data *dao.MenuList, err error) {
	data, err = dao.Menu.GetMenuList(title, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetUserMenu 获取用户菜单
func (m *menu) GetUserMenu(username string) (data []*dao.MenuItem, err error) {
	data, err = dao.Menu.GetUserMenu(username)
	if err != nil {
		return nil, err
	}

	return data, nil
}
