package dao

import (
	"errors"
	"gorm.io/gorm"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils"
)

var CasBin casbin

type casbin struct{}

// UpdateRolePermission 更新角色关联的权限（使用CasBin API操作，不支持事务，如果要支持事务需要改成原生Gorm操作）
func (c *casbin) UpdateRolePermission(groupName string, menus, paths []string) (err error) {

	// 迁移角色对应的所有权限
	_, err = global.CasBinServer.RemoveFilteredPolicy(0, groupName, "", "")
	if err != nil {
		return err
	}

	// 添加角色对应的接口权限
	for _, permission := range paths {

		// 获取前端传入的接口信息
		pathInfo, err := Path.GetPathInfo(permission)
		if err != nil {
			return err
		}

		// 添加角色对应的接口权限
		_, err = global.CasBinServer.AddNamedPolicy("p", groupName, pathInfo.Path, pathInfo.Method)
		if err != nil {
			return err
		}
	}

	// 添加角色对应的菜单权限
	for _, permission := range menus {

		_, err = global.CasBinServer.AddNamedPolicy("p", groupName, permission, "read")
		if err != nil {
			return err
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
