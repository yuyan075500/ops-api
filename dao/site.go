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

// SiteList 返回给站点列表结构体（表格）
type SiteList struct {
	Items []*SiteGroup `json:"items"`
	Total int64        `json:"total"`
}

// SiteGuideList 返回给站点列表结构体（站点导航）
type SiteGuideList struct {
	Items []*SiteGuideGroup `json:"items"`
}

// SiteGroup 站点分组（表格）
type SiteGroup struct {
	ID    uint        `json:"id"`
	Name  string      `json:"name"`
	Sites []*SiteItem `json:"sites"`
}

// SiteGuideGroup 站点分组（站点导航）
type SiteGuideGroup struct {
	ID    uint             `json:"id"`
	Name  string           `json:"name"`
	Sites []*SiteGuideItem `json:"sites"`
}

// SiteItem 站点（表格）
type SiteItem struct {
	ID           uint             `json:"id"`
	Name         string           `json:"name"`
	Icon         string           `json:"icon"`
	Address      string           `json:"address"`
	AllOpen      bool             `json:"all_open"`
	Description  string           `json:"description"`
	SSO          bool             `json:"sso"`
	SSOType      uint             `json:"sso_type"`
	ClientId     string           `json:"client_id"`
	ClientSecret string           `json:"client_secret"`
	CallbackUrl  string           `json:"callback_url"`
	EntityId     string           `json:"entity_id"`
	Certificate  string           `json:"certificate"`
	DomainId     string           `json:"domain_id"`
	RedirectUrl  string           `json:"redirect_url"`
	IDPName      string           `json:"idp_name"`
	Users        []*UserBasicInfo `json:"users"`
}

// SiteGuideItem 站点（站点导航）
type SiteGuideItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

// UpdateSite 更新站点结构体，定义新增时的字段信息
type UpdateSite struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	SSO         *bool   `json:"sso"`      // 指针类型，可以确保使用Updates方法更新时，如果值为false时也能更新成功
	AllOpen     *bool   `json:"all_open"` // 指针类型，可以确保使用Updates方法更新时，如果值为false时也能更新成功
	SSOType     uint    `json:"sso_type"`
	Icon        string  `json:"icon"`
	EntityId    *string `json:"entity_id"` // 指针类型，可以确保使用Updates方法更新时，如果值为空时也能更新成功
	CallbackUrl *string `json:"callback_url"`
	Certificate *string `json:"certificate"`
	Description string  `json:"description"`
}

// GetSiteGuideList 获取站点列表（站点导航）
func (s *site) GetSiteGuideList() (data *SiteGuideList, err error) {
	// 定义返回的内容
	var siteGroups []*model.SiteGroup

	// 获取分组列表
	tx := global.MySQLClient.Model(&model.SiteGroup{}).
		Preload("Sites"). // 预加载分组包含的站点
		Find(&siteGroups)
	if tx.Error != nil {
		return nil, errors.New(tx.Error.Error())
	}

	// 最外层结构体数据绑定（由于需要对站点URL特殊处理，所以不能直接返回siteGroups结果）
	siteList := &SiteGuideList{
		Items: make([]*SiteGuideGroup, len(siteGroups)), // 初始化分组列表切片，指定长度为siteGroups
	}

	// 对分组进行循环处理
	for i, sg := range siteGroups {
		siteGroup := &SiteGuideGroup{
			ID:    sg.ID,
			Name:  sg.Name,
			Sites: make([]*SiteGuideItem, len(sg.Sites)), // 初始化分组内的站点列表切片，指定长度为sg.Sites
		}

		// 对分组内的站点循环处理
		for j, s := range sg.Sites {
			siteItem := &SiteGuideItem{
				ID:          s.ID,
				Name:        s.Name,
				Address:     s.Address,
				Description: s.Description,
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

// GetSiteList 获取站点列表（表格）
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
		Preload("Sites").                   // 预加载分组包含的站点
		Preload("Sites.Users").             // 确保预加载站点用户
		Where("name like ?", "%"+name+"%"). // 实现过滤
		Count(&total).                      // 获取总数
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
				Address:      s.Address,
				AllOpen:      s.AllOpen,
				Description:  s.Description,
				SSO:          s.SSO,
				SSOType:      s.SSOType,
				ClientId:     s.ClientId,
				ClientSecret: s.ClientSecret,
				CallbackUrl:  s.CallbackUrl,
				EntityId:     s.EntityId,
				Certificate:  s.Certificate,
				DomainId:     s.DomainId,
				RedirectUrl:  s.RedirectUrl,
				IDPName:      s.IDPName,
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

			// 处理用户信息
			siteItem.Users = make([]*UserBasicInfo, len(s.Users))
			for k, u := range s.Users {
				siteItem.Users[k] = &UserBasicInfo{
					ID:   u.ID,
					Name: u.Name,
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

// AddGroup 新增站点分组
func (s *site) AddGroup(data *model.SiteGroup) (err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// AddSite 新增站点
func (s *site) AddSite(data *model.Site) (err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// UpdateGroup 修改站点分组
func (s *site) UpdateGroup(data *model.SiteGroup) (err error) {
	if err := global.MySQLClient.Model(&model.SiteGroup{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// UpdateSite 修改站点
func (s *site) UpdateSite(site *model.Site, data *UpdateSite) (err error) {
	if err := global.MySQLClient.Model(&site).Updates(data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// DeleteGroup 删除站点分组
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

// DeleteSite 删除站点
func (s *site) DeleteSite(site *model.Site) (err error) {

	// 开启事务
	tx := global.MySQLClient.Begin()

	// 删除站点内所有用户
	if err := tx.Model(&site).Association("Users").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// 删除站点
	if err := tx.Unscoped().Delete(&site).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// GetCASSite 获取单个使用CAS3.0认证的站点
func (s *site) GetCASSite(service string) (data *model.Site, err error) {
	var site *model.Site

	if err := global.MySQLClient.Where("callback_url = ? AND sso = true AND sso_type = 1", service).First(&site).Error; err != nil {
		return nil, err
	}

	return site, nil
}

// GetOAuthSite 获取单个使用OAuth2.0认证的站点
func (s *site) GetOAuthSite(clientId string) (data *model.Site, err error) {
	var site *model.Site

	if err := global.MySQLClient.Where("client_id = ? AND sso = true AND sso_type = 2", clientId).First(&site).Error; err != nil {
		return nil, err
	}

	return site, nil
}

// GetSamlSite 获取单个使用SAML2认证的站点
func (s *site) GetSamlSite(issuer string) (data *model.Site, err error) {
	var site *model.Site

	if err := global.MySQLClient.Where("entity_id = ? AND sso = true AND sso_type = 3", issuer).First(&site).Error; err != nil {
		return nil, err
	}

	return site, nil
}

// UpdateSiteUser 更新站点用户
func (s *site) UpdateSiteUser(site *model.Site, users []model.AuthUser) (err error) {
	if err := global.MySQLClient.Model(&site).Association("Users").Replace(users); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// ClearSiteUser 清空站点用户
func (s *site) ClearSiteUser(site *model.Site) (err error) {
	if err := global.MySQLClient.Model(&site).Association("Users").Clear(); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// IsUserInSite 判断用户是否在站点中
func (s *site) IsUserInSite(userID uint, site *model.Site) bool {

	// 查询站点并预加载用户
	if err := global.MySQLClient.Preload("Users", "id = ?", userID).First(&site).Error; err != nil {
		return false
	}

	// 检查用户是否被预加载
	for _, user := range site.Users {
		if user.ID == userID {
			return true
		}
	}

	return false
}
