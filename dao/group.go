package dao

import (
	"errors"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
	"ops-api/global"
	"ops-api/model"
)

var Group group

type group struct{}

// GroupList 返回给前端的结构体
type GroupList struct {
	Items []*GroupInfo `json:"items"`
	Total int64        `json:"total"`
}

// GroupInfo 返回字段信息
type GroupInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IsRoleGroup bool   `json:"is_role_group"`
}

// GetGroupList 获取列表
func (u *group) GetGroupList(name string, page, limit int) (data *GroupList, err error) {
	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		groupInfo []*GroupInfo
		total     int64
	)

	// 获取分组列表
	tx := global.MySQLClient.Model(&model.AuthGroup{}).
		Where("name like ?", "%"+name+"%"). // 实现过滤
		Count(&total).                      // 获取总数
		Limit(limit).
		Offset(startSet).
		Find(&groupInfo)
	if tx.Error != nil {
		return nil, errors.New(tx.Error.Error())
	}

	return &GroupList{
		Items: groupInfo,
		Total: total,
	}, nil
}

// AddGroup 新增
func (u *group) AddGroup(tx *gorm.DB, data *model.AuthGroup) (err error) {
	if err := tx.Create(&data).Error; err != nil {
		logger.Error("ERROR：", err.Error())
		return errors.New(err.Error())
	}
	return nil
}

// UpdateGroup 修改
func (u *group) UpdateGroup(tx *gorm.DB, data *model.AuthGroup) (err error) {
	if err := tx.Model(&model.AuthGroup{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		logger.Error("ERROR：", err.Error())
		return errors.New(err.Error())
	}
	return nil
}

// DeleteGroup 删除
func (u *group) DeleteGroup(tx *gorm.DB, group *model.AuthGroup) (err error) {

	// 清除关联关系
	if err := Group.ClearGroupUser(tx, group); err != nil {
		logger.Error("ERROR：", err)
		return err
	}

	// 删除分组
	if err := tx.Unscoped().Delete(&group).Error; err != nil {
		logger.Error("ERROR：", err.Error())
		return errors.New(err.Error())
	}
	return nil
}

// UpdateGroupUser 更新组用户
func (u *group) UpdateGroupUser(tx *gorm.DB, group *model.AuthGroup, users []model.AuthUser) (err error) {
	if err := tx.Model(&group).Association("Users").Replace(users); err != nil {
		logger.Error("ERROR：", err.Error())
		return errors.New(err.Error())
	}

	return nil
}

// ClearGroupUser 清空组用户
func (u *group) ClearGroupUser(tx *gorm.DB, group *model.AuthGroup) (err error) {
	if err := tx.Model(&group).Association("Users").Clear(); err != nil {
		logger.Error("ERROR：", err.Error())
		return errors.New(err.Error())
	}

	return nil
}
