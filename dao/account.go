package dao

import (
	"ops-api/global"
	"ops-api/model"
)

var Account account

type account struct{}

// AccountList 返回给前端表格的数据结构体
type AccountList struct {
	Items []*AccountInfo `json:"items"`
	Total int64          `json:"total"`
}
type AccountInfo struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	LoginAddress string       `json:"login_address"`
	LoginMethod  string       `json:"login_method"`
	Username     string       `json:"username"`
	Note         string       `json:"note"`
	AuthUserID   uint         `json:"aut h_user_id"`
	AccountOwner AccountOwner `json:"account_owner" gorm:"-"`
}
type AccountOwner struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// AddAccount 新增
func (a *account) AddAccount(data *model.Account) (err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

// GetAccountList 获取账号列表
func (a *account) GetAccountList(name string, userID uint, page, limit int) (data *AccountList, err error) {
	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		accountList []*AccountInfo
		total       int64
	)

	// 获取账号列表（只返回用户自己的和别人分享的）
	if err := global.MySQLClient.Model(&model.Account{}).
		Where("auth_user_id = ? AND (name like ? OR username like ? OR login_address like ? OR note like ?)", userID, "%"+name+"%", "%"+name+"%", "%"+name+"%", "%"+name+"%"). // 实现过滤
		Count(&total).                                                                                                                                                         // 获取总数
		Limit(limit).
		Offset(startSet).
		Find(&accountList).Error; err != nil {
		return nil, err
	}

	// 获取用户信息
	for _, account := range accountList {
		conditions := map[string]interface{}{
			"id": account.AuthUserID,
		}
		user, err := User.GetUser(conditions)
		if err != nil {
			return nil, err
		}
		account.AccountOwner = AccountOwner{
			ID:   user.ID,
			Name: user.Name,
		}
	}

	return &AccountList{
		Items: accountList,
		Total: total,
	}, nil
}

// GetAccountPassword 获取账号密码
func (a *account) GetAccountPassword(id uint) (password string, err error) {

	var str string

	// 查找账号密码
	if err := global.MySQLClient.Model(&model.Account{}).
		Select("password").
		Where("id = ?", id).
		Scan(&str).Error; err != nil {
		return "", err
	}

	return str, nil
}
