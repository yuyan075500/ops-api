package dao

import (
	"errors"
	"github.com/wonderivan/logger"
	"ops-api/global"
	"ops-api/model"
)

var Group group

type group struct{}

// GroupList 返回给前端的结构体
type GroupList struct {
	Items []*model.AuthGroup `json:"items"`
	Total int64              `json:"total"`
}

// GetGroupList 获取列表
func (u *group) GetGroupList(name string, page, limit int) (data *GroupList, err error) {
	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		groupList []*model.AuthGroup
		total     int64
	)

	// 获取用户列表
	tx := global.MySQLClient.Model(&model.AuthGroup{}).
		Where("name like ?", "%"+name+"%"). // 实现过滤
		Count(&total).                      // 获取总数
		Limit(limit).
		Offset(startSet).
		Find(&groupList)
	if tx.Error != nil {
		logger.Error("获取列表失败：", tx.Error)
		return nil, errors.New("获取列表失败：" + tx.Error.Error())
	}

	return &GroupList{
		Items: groupList,
		Total: total,
	}, nil
}

// AddGroup 新增
func (u *group) AddGroup(data *model.AuthGroup) (err error) {
	tx := global.MySQLClient.Create(&data)
	if tx.Error != nil {
		logger.Error("新增失败：", tx.Error)
		return errors.New("新增失败：" + tx.Error.Error())
	}
	return nil
}

// UpdateGroup 修改
func (u *group) UpdateGroup(data *model.AuthGroup) (err error) {
	tx := global.MySQLClient.Model(&model.AuthGroup{}).Where("id = ?", data.ID).Updates(data)
	if tx.Error != nil {
		logger.Error("更新失败：", tx.Error)
		return errors.New("更新失败：" + tx.Error.Error())
	}
	return nil
}

// DeleteGroup 删除
func (u *group) DeleteGroup(id int) (err error) {
	tx := global.MySQLClient.Where("id = ?", id).Unscoped().Delete(&model.AuthGroup{})
	if tx.Error != nil {
		logger.Error("删除失败：", tx.Error)
		return errors.New("删除失败：" + tx.Error.Error())
	}
	return nil
}

// UpdateGroupUser 更新组用户
func (u *group) UpdateGroupUser(group *model.AuthGroup, users []model.AuthUser) (err error) {
	if err := global.MySQLClient.Debug().Model(&group).Association("Users").Replace(users); err != nil {
		logger.Error("更新失败：", err.Error)
		return errors.New("更新失败：" + err.Error())
	}

	return nil
}
