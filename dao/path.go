package dao

import (
	"errors"
	"ops-api/global"
	"ops-api/model"
)

var Path path

type path struct{}

// PathList 返回给前端表格的数据结构体
type PathList struct {
	Items []*model.SystemPath `json:"items"`
	Total int64               `json:"total"`
}

type MenuPaths struct {
	MenuName string              `json:"menu_name"`
	Paths    []*model.SystemPath `json:"paths"`
}

// GetPathList 获取接口列表（表格）
func (p *path) GetPathList(menuName string, page, limit int) (data *PathList, err error) {
	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		pathList []*model.SystemPath
		total    int64
	)

	// 获取用户列表
	tx := global.MySQLClient.Model(&model.SystemPath{}).
		Where("menu_name = ?", menuName).
		Count(&total). // 获取总数
		Limit(limit).
		Offset(startSet).
		Find(&pathList)
	if tx.Error != nil {
		return nil, errors.New(tx.Error.Error())
	}

	return &PathList{
		Items: pathList,
		Total: total,
	}, nil
}

// GetPathListAll 获取接口列表（复选框）
func (p *path) GetPathListAll() (data []MenuPaths, err error) {
	var (
		paths     []*model.SystemPath
		menuPaths []MenuPaths
	)

	// 获取用户列表
	if err := global.MySQLClient.Model(&model.SystemPath{}).
		Find(&paths).Error; err != nil {
		return nil, errors.New(err.Error())
	}

	// 按名称分类
	result := make(map[string][]*model.SystemPath)
	for _, path := range paths {
		result[path.MenuName] = append(result[path.MenuName], path)
	}

	// 构建返回结果
	for menuName, pathList := range result {

		// 根据菜单Name获取对应的Title
		title, err := Menu.GetMenuTitle(menuName)
		if err != nil {
			return nil, err
		}

		path := MenuPaths{
			MenuName: *title,
			Paths:    pathList,
		}
		menuPaths = append(menuPaths, path)
	}

	return menuPaths, nil
}

// GetPathInfo 根据接口Name获取详情
func (p *path) GetPathInfo(name string) (data *model.SystemPath, err error) {

	var path *model.SystemPath

	tx := global.MySQLClient.Where("name = ?", name).First(&path)
	if tx.Error != nil {
		return nil, errors.New(tx.Error.Error())
	}

	return path, nil
}
