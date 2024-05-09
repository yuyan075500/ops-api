package service

import "ops-api/dao"

var Menu menu

type menu struct{}

// GetUserMenu 获取用户菜单
func (m *menu) GetUserMenu() (data []*dao.MenuItem, err error) {
	data, err = dao.Menu.GetUserMenu()
	if err != nil {
		return nil, err
	}

	return data, nil
}
