package dao

import (
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
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	IsRoleGroup bool             `json:"is_role_group"`
	Users       []*UserBasicInfo `json:"users"`
	Menus       []string         `json:"menus"`
	Paths       []string         `json:"paths"`
}

// GetGroupList 获取列表
func (u *group) GetGroupList(name string, page, limit int) (data *GroupList, err error) {
	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		authGroup []*model.AuthGroup
		total     int64
	)

	// 获取分组列表
	if err := global.MySQLClient.Model(&model.AuthGroup{}).
		Preload("Users").                   // 预加载用户信息
		Where("name like ?", "%"+name+"%"). // 实现过滤
		Count(&total).                      // 获取总数
		Limit(limit).
		Offset(startSet).
		Find(&authGroup).Error; err != nil {
		return nil, err
	}

	// 绑定最外层结构体的数据
	groupList := &GroupList{
		Total: total,
		Items: make([]*AuthGroup, len(authGroup)), // 初始化分组列表切片，并指定长度为authGroup长度
	}

	for g, group := range authGroup {

		// 获取分组的权限列表
		permissions := global.CasBinServer.GetFilteredPolicy(0, group.Name)
		menus := make([]string, 0)
		paths := make([]string, 0)
		for _, permission := range permissions {
			if permission[2] == "read" {
				v1 := permission[1]
				menus = append(menus, v1)
			} else {
				name, err := Path.GetPathName(permission[1], permission[2])
				if err != nil {
					return nil, err
				}
				paths = append(paths, *name)
			}
		}

		// 绑定分组相关数据到结构体
		groupItem := &AuthGroup{
			ID:          group.ID,
			Name:        group.Name,
			IsRoleGroup: group.IsRoleGroup,
			Users:       make([]*UserBasicInfo, len(group.Users)), // 初始化用户列表切片，并指定长度为group.Users长度
			Menus:       menus,
			Paths:       paths,
		}

		// 遍历用户列表，绑定用户数据到对应分组结构体
		for u, user := range group.Users {
			groupItem.Users[u] = &UserBasicInfo{
				ID:   user.ID,
				Name: user.Name,
			}
		}

		// 追加分组数据到列表
		groupList.Items[g] = groupItem
	}

	return groupList, nil
}

// AddGroup 新增
func (u *group) AddGroup(data *model.AuthGroup) (authGroup *model.AuthGroup, err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateGroup 修改
func (u *group) UpdateGroup(tx *gorm.DB, data *model.AuthGroup) (*model.AuthGroup, error) {
	if err := tx.Model(&model.AuthGroup{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}

	// 获取最新的数据库记录，确保UpdatedAt是最新的
	var updatedGroup model.AuthGroup
	if err := tx.First(&updatedGroup, data.ID).Error; err != nil {
		return nil, err
	}

	// 返回更新后的数据
	return &updatedGroup, nil
}

// DeleteGroup 删除
func (u *group) DeleteGroup(tx *gorm.DB, group *model.AuthGroup) (err error) {

	// 清除关联关系
	if err := Group.ClearGroupUser(tx, group); err != nil {
		return err
	}

	// 删除分组
	if err := tx.Unscoped().Delete(&group).Error; err != nil {
		return err
	}
	return nil
}

// UpdateGroupUser 更新组用户
func (u *group) UpdateGroupUser(tx *gorm.DB, group *model.AuthGroup, users []model.AuthUser) (*model.AuthGroup, error) {
	if err := tx.Model(&group).Association("Users").Replace(users); err != nil {
		return nil, err
	}
	return group, nil
}

// ClearGroupUser 清空组用户
func (u *group) ClearGroupUser(tx *gorm.DB, group *model.AuthGroup) (err error) {
	return tx.Model(&group).Association("Users").Clear()
}
