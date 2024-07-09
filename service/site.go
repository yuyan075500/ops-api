package service

import (
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
)

var Site site

type site struct{}

// SiteGroupCreate 创建站点分组结构体，定义新增时的字段信息
type SiteGroupCreate struct {
	Name string `json:"name" binding:"required"`
}

// SiteCreate 创建站点结构体，定义新增时的字段信息
type SiteCreate struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	SSO         *bool  `json:"sso" binding:"required"`
	SSOType     uint   `json:"sso_type"`
	Icon        string `json:"icon"`
	CallbackUrl string `json:"callback_url"`
	Description string `json:"description" binding:"required"`
	SiteGroupID uint   `json:"site_group_id" binding:"required"`
}

// SiteGroupUpdate 更新分组名称构体
type SiteGroupUpdate struct {
	ID   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// GetSiteList 获取站点分组列表（表格）
func (s *site) GetSiteList(name string, page, limit int) (data *dao.SiteList, err error) {
	data, err = dao.Site.GetSiteList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetSiteGuideList 获取站点分组列表（站点导航）
func (s *site) GetSiteGuideList() (data *dao.SiteGuideList, err error) {
	data, err = dao.Site.GetSiteGuideList()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// AddGroup 创建站点分组
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

// AddSite 创建站点
func (s *site) AddSite(data *SiteCreate) (err error) {

	group := &model.Site{
		Name:        data.Name,
		Address:     data.Address,
		SSO:         *data.SSO,
		SSOType:     data.SSOType,
		Icon:        &data.Icon,
		CallbackUrl: data.CallbackUrl,
		Description: data.Description,
		SiteGroupID: data.SiteGroupID,
	}

	// 创建数据库数据
	if err := dao.Site.AddSite(group); err != nil {
		return err
	}

	return nil
}

// UpdateGroup 更新站点分组
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

// UpdateSite 更新站点
func (s *site) UpdateSite(data *dao.UpdateSite) error {

	// 查询要修改的站点
	site := &model.Site{}
	if err := global.MySQLClient.First(site, data.ID).Error; err != nil {
		return err
	}

	return dao.Site.UpdateSite(site, data)
}

// DeleteGroup 删除站点分组
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

// DeleteSite 删除站点
func (s *site) DeleteSite(id int) (err error) {

	site := &model.Site{}
	if err := global.MySQLClient.First(site, id).Error; err != nil {
		return err
	}

	// 删除分组
	if err := dao.Site.DeleteSite(site); err != nil {
		return err
	}

	return nil
}
