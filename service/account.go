package service

import (
	"errors"
	"fmt"
	"github.com/pquerna/otp/totp"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils"
	message "ops-api/utils/sms"
	"time"
)

var Account account

type account struct{}

// AccountCreate 账号创建结构体
type AccountCreate struct {
	Name         string `json:"name" binding:"required"`
	LoginAddress string `json:"login_address"`
	LoginMethod  string `json:"login_method"`
	Username     string `json:"username"`
	Password     string `json:"password" binding:"required"`
	Note         string `json:"note"`
}

// CodeVerification 获取密码结构体
type CodeVerification struct {
	ValidateType uint   `json:"validate_type" binding:"required"` // 验证类型：1：短信验证码，2：MFA验证码
	Code         string `json:"code" binding:"required"`
}

// CheckAccountOwnership 校验用户是否有权操作指定账号
func (a *account) CheckAccountOwnership(accountID int, userID uint) error {
	user, err := dao.Account.GetAccountOwner(accountID)
	if err != nil {
		return err
	}
	if user.ID != userID {
		return errors.New("此账号你无权操作")
	}
	return nil
}

// AddAccount 创建账号
func (a *account) AddAccount(data *AccountCreate, userId uint) (err error) {

	account := &model.Account{
		Name:         data.Name,
		Username:     data.Username,
		Password:     data.Password,
		LoginAddress: data.LoginAddress,
		Note:         data.Note,
		LoginMethod:  data.LoginMethod,
		OwnerUserID:  userId,
	}

	// 创建数据库数据
	err = dao.Account.AddAccount(account)
	if err != nil {
		return err
	}
	return nil
}

// GetAccountList 获取账号列表（表格）
func (a *account) GetAccountList(name string, userID uint, page, limit int) (data *dao.AccountList, err error) {
	data, err = dao.Account.GetAccountList(name, userID, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteAccount 删除账号
func (a *account) DeleteAccount(id, userId int) error {

	if err := a.CheckAccountOwnership(id, uint(userId)); err != nil {
		return err
	}

	return dao.Account.DeleteAccount(id)
}

// UpdateAccount 更新账号
func (a *account) UpdateAccount(data *dao.AccountUpdate, userId uint) error {

	if err := a.CheckAccountOwnership(int(data.ID), userId); err != nil {
		return err
	}

	return dao.Account.UpdateAccount(data)
}

// GetAccountPassword 获取账号密码
func (a *account) GetAccountPassword(id uint, username string, userId uint) (password *string, err error) {

	if err := a.CheckAccountOwnership(int(id), userId); err != nil {
		return nil, err
	}

	// 判断是否需要认证，Redis缓存中指定的Key是否存在，存在则不需要认证，否则需要认证
	var keyName = fmt.Sprintf("%s_get_account_password_enabled", username)
	val, err := global.RedisClient.Exists(keyName).Result()
	if err != nil {
		return nil, err
	}
	// 0 表示不存在
	if val == 0 {
		return nil, err
	}

	// 获取密码
	p, err := dao.Account.GetAccountPassword(id)
	if err != nil {
		return nil, err
	}

	// 密码解密
	pwd := p
	str, err := utils.Decrypt(pwd)
	if err != nil {
		return nil, err
	}

	return &str, nil
}

// GetSMSCode 发送获取账号密码验证码
func (a *account) GetSMSCode(userID uint) (err error) {

	// 查找用户
	conditions := map[string]interface{}{
		"id": userID,
	}
	user, err := dao.User.GetUser(conditions)
	if err != nil {
		return err
	}

	var keyName = fmt.Sprintf("%s_get_account_password_verification_code", user.Username)

	// 判断Redis缓存中指定的Key是否存在
	val, err := global.RedisClient.Exists(keyName).Result()
	if err != nil {
		return err
	}
	// 1 表示已存在
	if val == 1 {
		// 判断Key的有效期，如果Key的有效期大于4分钟，表示在1分钟内发送过验证码，提示用户请勿频繁发送校验码
		ttl, err := global.RedisClient.TTL(keyName).Result()
		if err != nil {
			return err
		}
		if ttl.Minutes() > 4 {
			return errors.New("请勿频繁发送校验码")
		}
	}

	// 发送短信验证码
	data := &message.SendData{
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Note:        "密码查询",
	}
	code, err := SMS.SMSSend(data)
	if err != nil {
		return err
	}

	// 将验证码写入Redis缓存，如果已存在则会更新Key的值并刷新TTL
	if err := global.RedisClient.Set(keyName, code, 5*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

// CodeVerification 校验验证码
func (a *account) CodeVerification(userId uint, data *CodeVerification) (err error) {

	// 查找用户
	conditions := map[string]interface{}{
		"id": userId,
	}
	user, err := dao.User.GetUser(conditions)
	if err != nil {
		return err
	}

	// 短信验证码校验
	if data.ValidateType == 1 {

		// 从缓存中获取验证码
		result, err := global.RedisClient.Get(fmt.Sprintf("%s_get_account_password_verification_code", user.Username)).Result()
		if err != nil {
			return err
		}
		// 判断是否正确
		if result != data.Code {
			return errors.New("校验码错误")
		}
	}

	// MFA验证码校验
	if data.ValidateType == 2 {

		// 获取Secret
		if user.MFACode == nil {
			return errors.New("您还未绑定MFA")
		}

		// 校验MFA
		valid := totp.Validate(data.Code, *user.MFACode)
		if !valid {
			return errors.New("校验码错误")
		}
	}

	// 写入允许用户获取密码的Redis缓存
	if err := global.RedisClient.Set(fmt.Sprintf("%s_get_account_password_enabled", user.Username), true, 10*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}
