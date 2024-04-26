package service

import (
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
)

var Group group

type group struct{}

// GroupCreate 创建构体，定义新增时的字段信息
type GroupCreate struct {
	Name string `json:"name" binding:"required"`
}

// GroupUpdate 更新构体，定义更新时的字段信息
type GroupUpdate struct {
	ID   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// GroupUpdateUser 更新组用户构体，定义更新时的字段信息
type GroupUpdateUser struct {
	ID    uint   `json:"id" binding:"required"`
	Users []uint `json:"users" binding:"required"`
}

// GetGroupList 获取列表
func (u *group) GetGroupList(name string, page, limit int) (data *dao.GroupList, err error) {
	data, err = dao.Group.GetGroupList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// AddGroup 创建
func (u *group) AddGroup(data *GroupCreate) (err error) {

	group := &model.AuthGroup{
		Name: data.Name,
	}

	// 创建数据库数据
	err = dao.Group.AddGroup(group)
	if err != nil {
		return err
	}
	return nil
}

// DeleteGroup 删除
func (u *group) DeleteGroup(id int) (err error) {
	err = dao.Group.DeleteGroup(id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateGroup 更新
func (u *group) UpdateGroup(data *GroupUpdate) error {

	// 查询要修改的数据
	group := &model.AuthGroup{}
	if err := global.MySQLClient.First(group, data.ID).Error; err != nil {
		return err
	}

	// 更新指定字段的值
	group.Name = data.Name

	return dao.Group.UpdateGroup(group)
}

// UpdateGroupUser 更新组用户
func (u *group) UpdateGroupUser(data *GroupUpdateUser) (err error) {

	// 查询要修改的用户组
	group := &model.AuthGroup{}
	if err := global.MySQLClient.First(group, data.ID).Error; err != nil {
		return err
	}

	// 查询出要更新的所有用户
	var users []model.AuthUser
	if err := global.MySQLClient.Find(&users, data.Users).Error; err != nil {
		return err
	}

	return dao.Group.UpdateGroupUser(group, users)
}
