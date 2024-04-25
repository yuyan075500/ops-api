package service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
	"ops-api/utils/check"
)

var User user

type user struct{}

// UserLogin 用户登录结构体
type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserCreate 创建用户结构体，定义新增用户时的字段信息
type UserCreate struct {
	Name        string `json:"name" binding:"required"`
	Username    string `json:"username" gorm:"unique" binding:"required"`
	Password    string `json:"password" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required" validate:"phone"`
	Email       string `json:"email" binding:"required" validate:"email"`
}

// UserPasswordUpdate 用户密码更新结构体
type UserPasswordUpdate struct {
	ID         uint   `json:"id" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required"`
}

// UserUpdate 用户更新构体，定义更新用户时的字段信息
type UserUpdate struct {
	ID          uint   `json:"id" binding:"required"`
	PhoneNumber string `json:"phone_number" validate:"phone"`
	Email       string `json:"email" validate:"email"`
	IsActive    bool   `json:"is_active"`
}

// GetUserList 获取用户列表
func (u *user) GetUserList(name string, page, limit int) (data *dao.UserList, err error) {
	data, err = dao.User.GetUserList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetUser 获取用户信息
func (u *user) GetUser(userid uint) (user *dao.UserInfo, err error) {
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
	}

	// 创建数据库数据
	err = dao.User.AddUser(user)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser 删除用户
func (u *user) DeleteUser(id int) (err error) {
	err = dao.User.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUser 更新用户信息
func (u *user) UpdateUser(data *UserUpdate) error {

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

	// 更新指定字段的值
	user.PhoneNumber = data.PhoneNumber
	user.Email = data.Email
	user.IsActive = &data.IsActive

	return dao.User.UpdateUser(data.ID, user)
}

// UpdateUserPassword 更新用户密码
func (u *user) UpdateUserPassword(data *UserPasswordUpdate) (err error) {

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

	// 更新密码
	user.Password = data.Password

	return dao.User.UpdateUserPassword(data.ID, user)
}
