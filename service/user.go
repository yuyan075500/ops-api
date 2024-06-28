package service

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/middleware"
	"ops-api/model"
	"ops-api/utils"
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

// RestPassword 重置密码时用户信息绑定的结构体
type RestPassword struct {
	Username    string `json:"username" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Password    string `json:"password" binding:"required"`
	RePassword  string `json:"re_password" binding:"required"`
}

// UserSync 用户同步结构体，用于AD用户同步
type UserSync struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	IsActive    bool   `json:"is_active"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
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
func (u *user) AddUser(data *dao.UserCreate) (err error) {

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
		IsActive:    true,
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

	// 如果是AD域账号则访问AD进行用户密码重置
	if err := AD.LDAPUserResetPassword(user.Username, data.Password); err != nil {
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

// UpdateSelfPassword 用户重置密码
func (u *user) UpdateSelfPassword(data *RestPassword) (err error) {

	var (
		user    model.AuthUser
		keyName = fmt.Sprintf("%s_rest_password_verification_code", data.Username)
	)

	// 获取用户信息
	tx := global.MySQLClient.First(&user, "username = ? AND phone_number = ?", data.Username, data.PhoneNumber)
	if tx.Error != nil {
		return errors.New("用户名或手机号错误")
	}

	// 验证码校验
	result, err := global.RedisClient.Get(keyName).Result()
	if err != nil {
		return err
	}
	if result != data.Code {
		return errors.New("校验码错误")
	}

	// 执行密码重置
	userInfo := &dao.UserPasswordUpdate{
		ID:         user.ID,
		Password:   data.Password,
		RePassword: data.RePassword,
	}
	if err := u.UpdateUserPassword(userInfo); err != nil {
		return err
	}

	return nil
}

// Login 用户登录
func (u *user) Login(params *UserLogin) (token string, redirect *string, err error) {

	var user model.AuthUser

	// 本地用户认证
	if !params.LDAP {
		// 根据用户名查询用户
		if err := global.MySQLClient.First(&user, "username = ? AND user_from = ?", params.Username, "本地").First(&user).Error; err != nil {
			return "", nil, errors.New("用户不存在")
		}
		if user.CheckPassword(params.Password) == false {
			return "", nil, errors.New("用户密码错误")
		}
	}

	// AD域用户认证
	if params.LDAP {
		// 根据用户名查询用户
		if err := global.MySQLClient.First(&user, "username = ? AND user_from = ?", params.Username, "AD域").First(&user).Error; err != nil {
			return "", nil, errors.New("用户不存在")
		}
		if _, err := AD.LDAPUserAuthentication(params.Username, params.Password); err != nil {
			return "", nil, errors.New("用户密码错误或系统错误")
		}
	}

	// 判断用户是否禁用
	if user.IsActive == false {
		return "", nil, errors.New("拒绝登录，请联系管理员")
	}

	// 判断系统是否启用MFA认证
	if config.Conf.MFA.Enable == true {
		// 生成一个32位长度的随机字符串，当用户进行MFA检验的时候判断用户否认证通过
		token := utils.GenerateRandomString(32)

		// 将token写入Redis缓存，设置有效期为2分钟
		if err := global.RedisClient.Set(token, user.Username, 2*time.Minute).Err(); err != nil {
			return "", nil, err
		}

		// 判断用户是否已经绑定MFA，为空则未绑定。redirect为前端跳转页面名称，如果用户未绑定MFA则跳转到MFA绑定页面，否则跳转MFA认证页面
		redirect := "MFA_AUTH"
		if user.MFACode == nil {
			redirect = "MFA_ENABLE"
			return token, &redirect, nil
		}

		return token, &redirect, nil
	}

	// 生成用户Token
	token, err = middleware.GenerateJWT(user.ID, user.Name, user.Username)
	if err != nil {
		return "", nil, err
	}

	// 记录用户最后登录时间
	global.MySQLClient.Model(&user).Where("id = ?", user.ID).Update("last_login_at", time.Now())

	return token, nil, nil
}

// UserSync AD域用户同步
func (u *user) UserSync() (err error) {
	// 获取AD域中的用户
	users, err := AD.LDAPUserSync()
	if err != nil {
		return err
	}

	// 创建用户
	var usersList []*model.AuthUser
	for _, user := range users {
		usersList = append(usersList, &model.AuthUser{
			Username:    user.Username,
			Name:        user.Name,
			Email:       user.Email,
			Password:    user.Password,
			IsActive:    user.IsActive,
			PhoneNumber: user.PhoneNumber,
			UserFrom:    user.UserFrom,
		})
	}
	if err := dao.User.SyncUsers(usersList); err != nil {
		return err
	}
	return nil
}
