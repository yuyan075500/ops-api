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
	Items []*AuthGroup `json:"items"`
	Total int64        `json:"total"`
}

// AuthGroup 返回分组的字段信息，这里结构体名称必须和实际模型名称保持一致
type AuthGroup struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	IsRoleGroup bool             `json:"is_role_group"`
	Users       []*UserBasicInfo `json:"users"`
	Menus       []*string        `json:"menus"`
}

// GetGroupList 获取列表
func (u *group) GetGroupList(name string, page, limit int) (data *GroupList, err error) {
	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		authGroup []*AuthGroup
		total     int64
	)

	// 获取分组列表
	tx := global.MySQLClient.Model(&model.AuthGroup{}).
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).                                 // 加载用户信息
		Where("name like ?", "%"+name+"%"). // 实现过滤
		Count(&total).                      // 获取总数
		Limit(limit).
		Offset(startSet).
		Find(&authGroup)
	if tx.Error != nil {
		return nil, errors.New(tx.Error.Error())
	}

	return &GroupList{
		Items: authGroup,
		Total: total,
	}, nil
}

// AddGroup 新增
func (u *group) AddGroup(data *model.AuthGroup) (err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// UpdateGroup 修改
func (u *group) UpdateGroup(tx *gorm.DB, data *model.AuthGroup) (err error) {
	if err := tx.Model(&model.AuthGroup{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// DeleteGroup 删除
func (u *group) DeleteGroup(tx *gorm.DB, group *model.AuthGroup) (err error) {

	// 清除关联关系
	if err := Group.ClearGroupUser(tx, group); err != nil {
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
		return errors.New(err.Error())
	}

	return nil
}

// ClearGroupUser 清空组用户
func (u *group) ClearGroupUser(tx *gorm.DB, group *model.AuthGroup) (err error) {
	if err := tx.Model(&group).Association("Users").Clear(); err != nil {
		return errors.New(err.Error())
	}

	return nil
}
