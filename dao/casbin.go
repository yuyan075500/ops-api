package dao

import (
	"errors"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils"
)

var CasBin casbin

type casbin struct{}

// UpdateRoleMenuPermission 更新角色关联的菜单权限
func (c *casbin) UpdateRoleMenuPermission(tx *gorm.DB, groupName string, permissions []string) (err error) {

	// 查询当前角色中所有权限列表
	var oldPermissions []string
	if err := tx.Model(&model.CasbinRule{}).Select("v1").Where("ptype = ? AND v0 = ? AND v2 = ?", "p", groupName, "read").Pluck("v1", &oldPermissions).Error; err != nil {
		logger.Error("ERROR：", err.Error())
		return errors.New(err.Error())
	}

	// 前端传入的权限列表中有，数据库中没有，则添加
	for _, permission := range permissions {
		if !utils.Contains(oldPermissions, permission) {
			if err := tx.Create(&model.CasbinRule{
				Ptype: "p",
				V0:    groupName,
				V1:    permission,
				V2:    "read",
			}).Error; err != nil {
				logger.Error("ERROR：", err.Error())
				return errors.New(err.Error())
			}
		}
	}

	// 前端传入的权限列表中没有，数据库中有，则删除
	for _, existingPermission := range oldPermissions {
		if !utils.Contains(permissions, existingPermission) {
			if err := tx.Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?", "p", groupName, existingPermission, "read").Delete(&model.CasbinRule{}).Error; err != nil {
				logger.Error("ERROR：", err.Error())
				return errors.New("ERROR：" + err.Error())
			}
		}
	}

	return nil
}

// UpdateRolePathPermission 更新角色关联的接口权限
func (c *casbin) UpdateRolePathPermission(tx *gorm.DB, groupName string, permissions []string) (err error) {

	// 查询当前角色中所有权限列表
	var oldPermissions []string
	if err := tx.Model(&model.CasbinRule{}).Select("v1").Where("ptype = ? AND v0 = ? AND v2 != ?", "p", groupName, "read").Pluck("v1", &oldPermissions).Error; err != nil {
		logger.Error("ERROR：", err.Error())
		return errors.New(err.Error())
	}

	// 前端传入的权限列表中有，数据库中没有，则添加
	for _, permission := range permissions {
		if !utils.Contains(oldPermissions, permission) {

			// 获取接口详情
			pathInfo, err := Path.GetPathInfo(permission)
			if err != nil {
				return err
			}

			if err := tx.Create(&model.CasbinRule{
				Ptype: "p",
				V0:    groupName,
				V1:    pathInfo.Path,
				V2:    pathInfo.Method,
			}).Error; err != nil {
				tx.Rollback()
				logger.Error("ERROR：", err.Error())
				return errors.New(err.Error())
			}
		}
	}

	// 前端传入的权限列表中没有，数据库中有，则删除
	for _, existingPermission := range oldPermissions {
		if !utils.Contains(permissions, existingPermission) {
			if err := tx.Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 != ?", "p", groupName, existingPermission, "read").Delete(&model.CasbinRule{}).Error; err != nil {
				tx.Rollback()
				logger.Error("ERROR：", err.Error())
				return errors.New("ERROR：" + err.Error())
			}
		}
	}

	return nil
}

// UpdateRoleUser 更新角色关联的用户
func (c *casbin) UpdateRoleUser(tx *gorm.DB, groupName string, users []string) (err error) {

	// 查询当前角色中所有用户列表
	var oldUsernames []string
	if err := tx.Model(&model.CasbinRule{}).Select("v0").Where("ptype = ? AND v1 = ?", "g", groupName).Pluck("v0", &oldUsernames).Error; err != nil {
		return errors.New(err.Error())
	}

	// 当前角色，前端传入的用户列表中有，数据库中没有，则添加
	for _, username := range users {
		if !utils.Contains(oldUsernames, username) {
			if err := tx.Create(&model.CasbinRule{
				Ptype: "g",
				V0:    username,
				V1:    groupName,
			}).Error; err != nil {
				return errors.New(err.Error())
			}
		}
	}

	// 当前角色，前端传入的用户列表中没有，数据库中有，则删除
	for _, existingUser := range oldUsernames {
		if !utils.Contains(users, existingUser) {
			if err := tx.Where("ptype = ? AND v0 = ? AND v1 = ?", "g", existingUser, groupName).Delete(&model.CasbinRule{}).Error; err != nil {
				return errors.New(err.Error())
			}
		}
	}

	// 加载规则
	if err := global.CasBinServer.LoadPolicy(); err != nil {
		return err
	}

	return nil
}

// UpdateRoleName 修改角色名称
func (c *casbin) UpdateRoleName(tx *gorm.DB, oldName, newName string) (err error) {

	if err := tx.Model(&model.CasbinRule{}).Where("ptype = ? AND v1 = ?", "g", oldName).Update("v1", newName).Error; err != nil {
		return errors.New(err.Error())
	}

	if err := tx.Model(&model.CasbinRule{}).Where("ptype = ? AND v0 = ?", "p", oldName).Update("v0", newName).Error; err != nil {
		return errors.New(err.Error())
	}

	// 加载规则
	if err := global.CasBinServer.LoadPolicy(); err != nil {
		return err
	}

	return nil
}

// DeleteRole 删除，删除所有与角色相关的记录
func (c *casbin) DeleteRole(tx *gorm.DB, groupName string) (err error) {

	if err := tx.Where("v0 = ? OR v1 = ?", groupName, groupName).Delete(&model.CasbinRule{}).Error; err != nil {
		return errors.New(err.Error())
	}

	// 加载规则
	if err := global.CasBinServer.LoadPolicy(); err != nil {
		return err
	}

	return nil
}
