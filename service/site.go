package service

import (
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
)

var Site site

type site struct{}

// SiteGroupCreate 创建构体，定义新增时的字段信息
type SiteGroupCreate struct {
	Name string `json:"name" binding:"required"`
}

// SiteGroupUpdate 更新分组名称构体
type SiteGroupUpdate struct {
	ID   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// GetSiteList 获取站点分组列表
func (s *site) GetSiteList(name string, page, limit int) (data *dao.SiteList, err error) {
	data, err = dao.Site.GetSiteList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// AddGroup 创建分组
func (s *site) AddGroup(data *SiteGroupCreate) (err error) {

	group := &model.SiteGroup{
		Name: data.Name,
	}

	// 创建数据库数据
	if err := dao.Site.AddGroup(group); err != nil {
		return err
	}

	return nil
}

// UpdateGroup 更新
func (s *site) UpdateGroup(data *SiteGroupUpdate) error {

	// 查询要修改的分组
	group := &model.SiteGroup{}
	if err := global.MySQLClient.First(group, data.ID).Error; err != nil {
		return err
	}

	// 更新分组名称
	group.Name = data.Name
	if err := dao.Site.UpdateGroup(group); err != nil {
		return err
	}

	return nil
}

func (s *site) DeleteGroup(id int) (err error) {

	group := &model.SiteGroup{}
	if err := global.MySQLClient.First(group, id).Error; err != nil {
		return err
	}

	// 删除分组
	if err := dao.Site.DeleteGroup(group); err != nil {
		return err
	}

	return nil
}
