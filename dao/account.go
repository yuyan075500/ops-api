package dao

import (
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils"
)

var Account account

type account struct{}

// AccountList 返回给前端表格的数据结构体
type AccountList struct {
	Items []*AccountInfo `json:"items"`
	Total int64          `json:"total"`
}
type AccountInfo struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	LoginAddress string          `json:"login_address"`
	LoginMethod  string          `json:"login_method"`
	Username     string          `json:"username"`
	Note         string          `json:"note"`
	OwnerUserID  uint            `json:"owner_user_id"`
	Owner        Owner           `json:"owner" gorm:"-"`
	Users        []*AccountUsers `json:"users"`
}
type Owner struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
type AccountUsers struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// AccountUpdate 更新账号信息
type AccountUpdate struct {
	ID           uint   `json:"id" binding:"required"`
	Name         string `json:"name"`
	LoginAddress string `json:"login_address"`
	LoginMethod  string `json:"login_method"`
	Username     string `json:"username"`
	OwnerUserID  uint   `json:"owner_user_id"`
	Note         string `json:"note"`
}

// AccountUpdateUser 更新账号共享用户
type AccountUpdateUser struct {
	ID    uint   `json:"id" binding:"required"`
	Users []uint `json:"users" binding:"required"`
}

// AccountUpdatePassword 更改密码
type AccountUpdatePassword struct {
	ID         uint   `json:"id" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required"`
}

// AddAccount 新增账号
func (a *account) AddAccount(data *model.Account) (account *model.Account, err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// AddAccounts 批量新增账号
func (a *account) AddAccounts(accounts []model.Account) (account []model.Account, err error) {
	if err := global.MySQLClient.Create(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

// DeleteAccount 删除账号
func (a *account) DeleteAccount(id int) (err error) {
	return global.MySQLClient.Where("id = ?", id).Unscoped().Delete(&model.Account{}).Error
}

// UpdateAccount 修改账号
func (a *account) UpdateAccount(data *AccountUpdate) (err error) {
	return global.MySQLClient.Model(&model.Account{}).Select("*").Where("id = ?", data.ID).Updates(data).Error
}

// BatchUpdateAccountOwner 批量修改账号所有者
func (a *account) BatchUpdateAccountOwner(accounts []uint, ownerId uint) (err error) {

	// 开启事务
	tx := global.MySQLClient.Begin()

	// 执行批量更新操作
	for _, accountID := range accounts {
		if err := tx.Model(&model.Account{}).
			Where("id = ?", accountID).
			Update("owner_user_id", ownerId).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	return tx.Commit().Error
}

// UpdatePassword 更改密码
func (a *account) UpdatePassword(account *model.Account, data *AccountUpdatePassword) (err error) {

	// 对密码进行加密
	cipherText, err := utils.Encrypt(data.Password)
	if err != nil {
		return err
	}
	return global.MySQLClient.Model(&account).Update("password", cipherText).Error
}

// UpdateAccountUser 账号分享
func (a *account) UpdateAccountUser(account *model.Account, users []model.AuthUser) (err error) {

	// 用户为空，则清空关联
	if len(users) == 0 {
		return global.MySQLClient.Model(&account).Association("Users").Clear()
	}

	return global.MySQLClient.Model(&account).Association("Users").Replace(users)
}

// GetAccountList 获取账号列表
func (a *account) GetAccountList(name string, userID uint, page, limit int) (*AccountList, error) {
	startSet := (page - 1) * limit

	var (
		accounts []*model.Account
		total    int64
	)

	// 查询数据库，获取满足条件的 Account 列表
	if err := global.MySQLClient.Model(&model.Account{}).
		Preload("OwnerUser").
		Preload("Users").
		Where("name LIKE ? OR username LIKE ? OR login_address LIKE ? OR note LIKE ?", "%"+name+"%", "%"+name+"%", "%"+name+"%", "%"+name+"%").
		Where("owner_user_id = ? OR id IN (?)",
			userID,
			global.MySQLClient.Table("account_users").Select("account_id").Where("auth_user_id = ?", userID),
		).
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Find(&accounts).Error; err != nil {
		return nil, err
	}

	var accountList []*AccountInfo

	// 将Account数据二次处理，转换为AccountInfo
	for _, account := range accounts {
		// 构建Owner
		owner := Owner{
			ID:   account.OwnerUser.ID,
			Name: account.OwnerUser.Name,
		}

		// 构建Users

		var users []*AccountUsers
		for _, user := range account.Users {
			users = append(users, &AccountUsers{
				ID:   user.ID,
				Name: user.Name,
			})
		}

		// 添加转换后的AccountInfo到列表中
		accountList = append(accountList, &AccountInfo{
			ID:           int(account.ID),
			Name:         account.Name,
			LoginAddress: account.LoginAddress,
			LoginMethod:  account.LoginMethod,
			Username:     account.Username,
			Note:         account.Note,
			OwnerUserID:  account.OwnerUserID,
			Owner:        owner,
			Users:        users,
		})
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

// GetAccountOwnerAndUsers 查询账号所有者和用户
func (a *account) GetAccountOwnerAndUsers(id int) (owner *model.AuthUser, users []*model.AuthUser, err error) {
	var account model.Account

	if err := global.MySQLClient.Preload("OwnerUser").Preload("Users").Where("id = ?", id).First(&account).Error; err != nil {
		return nil, nil, err
	}

	return &account.OwnerUser, account.Users, nil
}
