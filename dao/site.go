package dao

import (
	"errors"
	"ops-api/config"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils"
	"time"
)

var Site site

type site struct{}

// SiteList 返回给站点列表结构体
type SiteList struct {
	Items []*SiteGroup `json:"items"`
	Total int64        `json:"total"`
}

// SiteGroup 站点分组
type SiteGroup struct {
	ID    uint        `json:"id"`
	Name  string      `json:"name"`
	Sites []*SiteItem `json:"sites"`
}

// SiteItem 站点
type SiteItem struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	Address      string `json:"address"`
	AllOpen      bool   `json:"all_open"`
	Description  string `json:"description"`
	SSO          bool   `json:"sso"`
	SSOType      string `json:"sso_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackUrl  string `json:"callback_url"`
}

func (s *site) GetSiteList(name string, page, limit int) (data *SiteList, err error) {
	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		siteGroups []*model.SiteGroup
		total      int64
	)

	// 获取分组列表
	tx := global.MySQLClient.Model(&model.SiteGroup{}).
		Preload("Sites"). // 预加载分组包含的站点
		Where("name like ?", "%"+name+"%"). // 实现过滤
		Count(&total). // 获取总数
		Limit(limit).
		Offset(startSet).
		Find(&siteGroups)
	if tx.Error != nil {
		return nil, errors.New(tx.Error.Error())
	}

	// 最外层结构体数据绑定（由于需要对站点URL特殊处理，所以不能直接返回siteGroups结果）
	siteList := &SiteList{
		Total: total,
		Items: make([]*SiteGroup, len(siteGroups)), // 初始化分组列表切片，指定长度为siteGroups
	}

	// 对分组进行循环处理
	for i, sg := range siteGroups {
		siteGroup := &SiteGroup{
			ID:    sg.ID,
			Name:  sg.Name,
			Sites: make([]*SiteItem, len(sg.Sites)), // 初始化分组内的站点列表切片，指定长度为sg.Sites
		}

		// 对分组内的站点循环处理
		for j, s := range sg.Sites {
			siteItem := &SiteItem{
				ID:           s.ID,
				Name:         s.Name,
				Icon:         *s.Icon,
				Address:      s.Address,
				AllOpen:      s.AllOpen,
				Description:  s.Description,
				SSO:          s.SSO,
				SSOType:      s.SSOType,
				ClientId:     s.ClientId,
				ClientSecret: s.ClientSecret,
				CallbackUrl:  s.CallbackUrl,
			}

			// 对站点图标进行特殊处理，返回一个Minio中的临时URL链接
			if s.Icon != nil {
				iconUrl, err := utils.GetPresignedURL(*s.Icon, time.Duration(config.Conf.JWT.Expires)*time.Hour)
				if err != nil {
					siteItem.Icon = ""
				} else {
					siteItem.Icon = iconUrl.String()
				}
			}

			// 追加站点到分组
			siteGroup.Sites[j] = siteItem
		}

		// 追加分组到返回给前端的结构体
		siteList.Items[i] = siteGroup
	}

	return siteList, nil
}

// AddGroup 新增
func (s *site) AddGroup(data *model.SiteGroup) (err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// UpdateGroup 修改
func (s *site) UpdateGroup(data *model.SiteGroup) (err error) {
	if err := global.MySQLClient.Model(&model.SiteGroup{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// DeleteGroup 删除
func (s *site) DeleteGroup(group *model.SiteGroup) (err error) {

	// 删除分组
	if err := global.MySQLClient.Unscoped().Delete(&group).Error; err != nil {

		// 如果分组中包含站点，则返回对应的提示信息
		if utils.IsForeignKeyConstraintError(err) {
			return errors.New("请确保分组中不包含站点")
		}

		return errors.New(err.Error())
	}
	return nil
}
