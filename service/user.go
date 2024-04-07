package service

import (
	"ops-api/dao"
	"ops-api/model"
)

var User user

type user struct{}

// UserCreate 创建用户的结构体
type UserCreate struct {
	Name        string `json:"name"`
	UserName    string `json:"userName" gorm:"unique"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
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
		UserName:    data.UserName,
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
