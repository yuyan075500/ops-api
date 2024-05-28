package service

import (
	"ops-api/dao"
)

var Path path

type path struct{}

// GetPathListAll 获取接口列表（复选框）
func (p *path) GetPathListAll() (data []dao.MenuPaths, err error) {
	data, err = dao.Path.GetPathListAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetPathList 获取接口列表（表格）
func (p *path) GetPathList(menuName string, page, limit int) (data *dao.PathList, err error) {
	data, err = dao.Path.GetPathList(menuName, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}
