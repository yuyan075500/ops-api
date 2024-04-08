package dao

import (
	"errors"
	"github.com/wonderivan/logger"
	"ops-api/db"
	"ops-api/model"
)

var User user

type user struct{}

type UserList struct {
	Items []*model.AuthUser `json:"items"`
	Total int64             `json:"total"`
}

// GetUserList 获取用户列表
func (u *user) GetUserList(name string, page, limit int) (data *UserList, err error) {
	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		userList []*model.AuthUser
		total    int64
	)

	// 获取用户列表
	tx := db.GORM.Model(&model.AuthUser{}).
		Where("name like ?", "%"+name+"%"). // 实现过滤
		Count(&total).                      // 获取总数
		Limit(limit).
		Offset(startSet).
		Find(&userList)
	if tx.Error != nil {
		logger.Error("获取用户列表失败：", tx.Error)
		return nil, errors.New("获取用户列表失败：" + tx.Error.Error())
	}

	return &UserList{
		Items: userList,
		Total: total,
	}, nil
}

// AddUser 新增用户
func (u *user) AddUser(data *model.AuthUser) (err error) {
	tx := db.GORM.Create(&data)
	if tx.Error != nil {
		logger.Error("新增用户失败：", tx.Error)
		return errors.New("新增用户失败：" + tx.Error.Error())
	}
	return nil
}
