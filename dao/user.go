package dao

import (
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"ops-api/config"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils"
	"time"
)

var User user

type user struct{}

// UserList 返回给前端表格的数据结构体
type UserList struct {
	Items []*UserInfo `json:"items"`
	Total int64       `json:"total"`
}

// UserInfo 用户信息结构体
type UserInfo struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Username    string     `json:"username"`
	PhoneNumber string     `json:"phone_number"`
	IsActive    bool       `json:"is_active"`
	Email       string     `json:"email"`
	Avatar      string     `json:"avatar"`
	LastLoginAt *time.Time `json:"last_login_at"`
	UserFrom    string     `json:"user_from"`
}

// UserListAll 返回给前端下拉框或穿梭框的数据结构体
type UserListAll struct {
	Users []*UserBasicInfo `json:"users"`
}

// UserBasicInfo 用户基本信息结构体
type UserBasicInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// UserUpdate 更新构体，定义更新时的字段信息
type UserUpdate struct {
	ID          uint   `json:"id" binding:"required"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,phone"`
	Email       string `json:"email" validate:"omitempty,email"`
	IsActive    bool   `json:"is_active" validate:"omitempty"`
}

// UserPasswordUpdate 更改密码结构体
type UserPasswordUpdate struct {
	ID         uint   `json:"id" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required"`
}

// GetUserListAll 获取所有用户
func (u *user) GetUserListAll() (data *UserListAll, err error) {

	var userBasicInfo []*UserBasicInfo

	// 获取用户列表
	if err := global.MySQLClient.Model(&model.AuthUser{}).
		Select("id, name").
		Find(&userBasicInfo).Error; err != nil {
		logger.Error("获取列表失败：", err.Error())
		return nil, errors.New(err.Error())
	}

	return &UserListAll{
		Users: userBasicInfo,
	}, nil
}

// GetUserList 获取用户列表
func (u *user) GetUserList(name string, page, limit int) (data *UserList, err error) {
	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		userList []*UserInfo
		total    int64
	)

	// 获取用户列表
	tx := global.MySQLClient.Model(&model.AuthUser{}).
		Where("name like ?", "%"+name+"%"). // 实现过滤
		Count(&total).                      // 获取总数
		Limit(limit).
		Offset(startSet).
		Find(&userList)
	if tx.Error != nil {
		logger.Error("获取列表失败：", tx.Error)
		return nil, errors.New("获取列表失败：" + tx.Error.Error())
	}

	return &UserList{
		Items: userList,
		Total: total,
	}, nil
}

// GetUser 获取用户信息
func (u *user) GetUser(userid uint) (user *UserInfo, err error) {

	var userInfo *UserInfo

	tx := global.MySQLClient.Model(&model.AuthUser{}).Where("id = ?", userid).Find(&userInfo)
	if tx.Error != nil {
		logger.Error("获取信息失败：", tx.Error)
		return nil, errors.New("获取信息失败：" + tx.Error.Error())
	}

	// 从OSS中获取头像临时访问URL，临时URL的过期时间与用户Token过期时间保持一致
	avatarURL, err := utils.GetPresignedURL(userInfo.Avatar, time.Duration(config.Conf.JWT.Expires)*time.Hour)
	if err != nil {
		logger.Error("获取用户头像失败：", err.Error())
		userInfo.Avatar = ""
		return userInfo, nil
	}

	userInfo.Avatar = avatarURL.String()
	return userInfo, nil
}

// AddUser 新增
func (u *user) AddUser(data *model.AuthUser) (err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// UpdateUser 修改
func (u *user) UpdateUser(user *model.AuthUser, data *UserUpdate) (err error) {
	fmt.Println(data)
	// 当is_active=0，需要使用Select选中对应字段进行更新，否则无法设置为0
	if err := global.MySQLClient.Model(&user).Select("*").Updates(data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// DeleteUser 删除
func (u *user) DeleteUser(id int) (err error) {
	if err := global.MySQLClient.Where("id = ?", id).Unscoped().Delete(&model.AuthUser{}).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// UpdateUserPassword 更改密码
func (u *user) UpdateUserPassword(user *model.AuthUser, data *UserPasswordUpdate) (err error) {

	// 对密码进行加密
	cipherText, err := utils.Encrypt(data.Password)
	if err != nil {
		return err
	}

	// 更新密码
	if err := global.MySQLClient.Model(&user).Update("password", cipherText).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// ResetUserMFA 重置MFA
func (u *user) ResetUserMFA(data *model.AuthUser) (err error) {

	// 将MFA重置为nil
	tx := global.MySQLClient.Model(&model.AuthUser{}).Where("id = ?", data.ID).Update("mfa_code", nil)
	if tx.Error != nil {
		logger.Error("重置失败：", tx.Error)
		return errors.New("重置失败：" + tx.Error.Error())
	}
	return nil
}
