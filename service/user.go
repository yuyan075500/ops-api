package service

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils/check"
	"strconv"
	"time"
)

var User user

type user struct{}

// UserLogin 用户登录结构体
type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	LDAP     bool   `json:"ldap"`
}

// UserCreate 创建结构体，定义新增时的字段信息
type UserCreate struct {
	Name        string `json:"name" binding:"required"`
	Username    string `json:"username" gorm:"unique" binding:"required"`
	Password    string `json:"password" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required" validate:"phone"`
	Email       string `json:"email" binding:"required" validate:"email"`
	UserFrom    string `json:"user_from"`
}

// GetUserListAll 获取用户列表（下拉框、穿梭框）
func (u *user) GetUserListAll() (data *dao.UserListAll, err error) {
	data, err = dao.User.GetUserListAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetUserList 获取用户列表（表格）
func (u *user) GetUserList(name string, page, limit int) (data *dao.UserList, err error) {
	data, err = dao.User.GetUserList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetUser 获取用户信息
func (u *user) GetUser(userid uint) (user *dao.UserInfoWithMenu, err error) {
	user, err = dao.User.GetUser(userid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// AddUser 创建用户
func (u *user) AddUser(data *UserCreate) (err error) {

	// 字段校验
	validate := validator.New()
	// 注册自定义检验方法
	if err := validate.RegisterValidation("phone", check.PhoneNumberCheck); err != nil {
		return err
	}
	if err := validate.Struct(data); err != nil {
		return err.(validator.ValidationErrors)
	}

	// 检查密码是否复合要求
	if err := check.PasswordCheck(data.Password); err != nil {
		return err
	}

	user := &model.AuthUser{
		Name:        data.Name,
		Username:    data.Username,
		Password:    data.Password,
		PhoneNumber: data.PhoneNumber,
		Email:       data.Email,
		UserFrom:    data.UserFrom,
	}

	// 创建数据库数据
	err = dao.User.AddUser(user)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser 删除
func (u *user) DeleteUser(id int) (err error) {
	err = dao.User.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUser 更新
func (u *user) UpdateUser(data *dao.UserUpdate) error {

	// 字段校验
	validate := validator.New()
	// 注册自定义检验方法
	if err := validate.RegisterValidation("phone", check.PhoneNumberCheck); err != nil {
		return err
	}
	if err := validate.Struct(data); err != nil {
		return err.(validator.ValidationErrors)
	}

	// 查询要修改的用户
	user := &model.AuthUser{}
	if err := global.MySQLClient.First(user, data.ID).Error; err != nil {
		return err
	}

	return dao.User.UpdateUser(user, data)
}

// UpdateUserPassword 更改密码
func (u *user) UpdateUserPassword(data *dao.UserPasswordUpdate) (err error) {

	// 检查密码校验
	if data.Password != data.RePassword {
		return errors.New("两次输入的密码不匹配")
	}
	if err := check.PasswordCheck(data.Password); err != nil {
		return err
	}

	// 查询要修改的用户
	user := &model.AuthUser{}
	if err := global.MySQLClient.First(user, data.ID).Error; err != nil {
		return err
	}

	return dao.User.UpdateUserPassword(user, data)
}

// ResetUserMFA 重置MFA
func (u *user) ResetUserMFA(id int) (err error) {

	// 查询要重置的用户
	user := &model.AuthUser{}
	if err := global.MySQLClient.First(user, id).Error; err != nil {
		return err
	}

	return dao.User.ResetUserMFA(user)
}

// GetVerificationCode 获取重置密码短信验证码
func (u *user) GetVerificationCode(data *UserInfo, expirationTime int) (err error) {

	var (
		keyName = fmt.Sprintf("%s_rest_password_verification_code", data.Username)
	)

	// 判断Redis缓存中指定的Key是否存在
	val, err := global.RedisClient.Exists(keyName).Result()
	if err != nil {
		return err
	}
	// 已存在
	if val >= 1 {
		// 判断Key的有效期
		ttl, err := global.RedisClient.TTL(keyName).Result()
		if err != nil {
			return err
		}
		// 如果Key的有效期大于4分钟，则提示用户请勿频繁发送校验码
		if ttl.Minutes() > 4 {
			return errors.New("请勿频繁发送校验码")
		}
	}

	// 发送短信验证码
	code, err := Log.SMSSend(data, strconv.Itoa(expirationTime))
	if err != nil {
		return err
	}

	// 将验证码写入Redis缓存，如果已存在则会更新Key的值并刷新TTL
	if err := global.RedisClient.Set(keyName, code, time.Duration(expirationTime)*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}
