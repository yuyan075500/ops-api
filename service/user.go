package service

import (
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
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
	PhoneNumber string `json:"phone_number" binding:"required"`
	Email       string `json:"email" binding:"required"`
}

// UserUpdate 用户更新构体，定义更新用户时的字段信息
type UserUpdate struct {
	ID          uint   `json:"id" binding:"required"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
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

// AddUser 创建用户
func (u *user) AddUser(data *UserCreate) (err error) {
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
func UpdateUser(data *UserUpdate) error {

	// 查询要修改的用户
	user := &model.AuthUser{}
	if err := global.MySQLClient.First(user, data.ID).Error; err != nil {
		return err
	}

	// 更新指定字段的值
	user.PhoneNumber = data.PhoneNumber
	user.Email = data.Email
	user.IsActive = data.IsActive

	return dao.User.UpdateUser(data.ID, user)
}
