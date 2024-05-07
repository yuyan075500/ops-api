package service

import (
	"errors"
	"gorm.io/gorm"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
)

var Group group

type group struct{}

// GroupCreate 创建构体，定义新增时的字段信息
type GroupCreate struct {
	Name        string `json:"name" binding:"required"`
	IsRoleGroup bool   `json:"is_role_group" default:"false"`
	Users       []uint `json:"users"`
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

// RoleCreate CasBin角色创建构体
type RoleCreate struct {
	Ptype string `json:"ptype"`
	V0    string `json:"v0"`
	V1    string `json:"v1"`
}

// GetGroupList 获取列表
func (u *group) GetGroupList(name string, page, limit int) (data *dao.GroupList, err error) {
	data, err = dao.Group.GetGroupList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// AddGroup 创建，支持同时添加用户
func (u *group) AddGroup(data *GroupCreate) (err error) {

	group := &model.AuthGroup{
		Name:        data.Name,
		IsRoleGroup: data.IsRoleGroup,
	}

	if data.IsRoleGroup && len(data.Users) == 0 {
		return errors.New("角色用户组必须至少添加一个用户")
	}

	// 如果传的的用户不为空，则添加用户
	if len(data.Users) > 0 {
		group.Users = make([]*model.AuthUser, len(data.Users))
		for index, userId := range data.Users {
			group.Users[index] = &model.AuthUser{
				Model: gorm.Model{
					ID: userId,
				},
			}
		}
	}

	// 开启事务
	tx := global.MySQLClient.Begin()

	// 创建数据库数据
	if err := dao.Group.AddGroup(tx, group); err != nil {
		tx.Rollback()
		return err
	}

	// 同步角色用户组信息到CasBin策略表
	users, err := GetUserNamesFromIDs(data.Users)
	if data.IsRoleGroup {
		for _, username := range users {
			rule := &model.CasbinRule{
				Ptype: "g",
				V0:    username,
				V1:    data.Name,
			}
			if err := dao.CasBin.AddRole(tx, rule); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
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

// GetUserNamesFromIDs 根据用户ID列表返回对应的用户名列表
func GetUserNamesFromIDs(userIDs []uint) ([]string, error) {
	var usernames []string

	err := global.MySQLClient.Model(&model.AuthUser{}).Select("username").Where("id IN (?)", userIDs).Find(&usernames).Error
	if err != nil {
		return nil, err
	}

	return usernames, nil
}
