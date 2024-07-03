package service

import "ops-api/dao"

var Site site

type site struct{}

// GetSiteList 获取站点分组列表
func (s *site) GetSiteList(name string, page, limit int) (data *dao.SiteList, err error) {
	data, err = dao.Site.GetSiteList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}
