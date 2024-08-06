package service

import (
	"errors"
	"fmt"
	"ops-api/dao"
	"ops-api/middleware"
	"ops-api/model"
	"ops-api/utils"
	"time"
)

var SSO sso

type sso struct{}

// Authorize 获取授权请求参数
type Authorize struct {
	ResponseType string `json:"response_type" binding:"required"`
	ClientId     string `json:"client_id" binding:"required"`
	RedirectURI  string `json:"redirect_uri"`
	State        string `json:"state"`
	Scope        string `json:"scope"`
}

// Token 获取token请求参数
type Token struct {
	GrantType    string `form:"grant_type"`
	Code         string `form:"code"`
	ClientId     string `form:"client_id"`
	RedirectURI  string `form:"redirect_uri"`
	ClientSecret string `form:"client_secret"`
}

// ResponseToken 返回给客户端的Token信息
type ResponseToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// ResponseUserinfo 返回给客户端的用户信息
type ResponseUserinfo struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

// GetAuthorize 客户端授权
func (s *sso) GetAuthorize(data *Authorize, userId uint) (callbackUrl string, err error) {

	// 获取客户端应用
	site, err := dao.Site.GetSite(data.ClientId)
	if err != nil {
		return "", err
	}

	// 判断用户是否有权限访问
	if !site.AllOpen {
		if !dao.Site.IsUserInSite(userId) {
			return "", errors.New("您无权访问该应用")
		}
	}

	// 创建随机字符串（长度建议>16）
	str := utils.GenerateRandomString(32)
	// 字符串加密，用于返回给客户端授权码
	code, err := utils.Encrypt(str)
	if err != nil {
		return "", err
	}

	// 将授权票据写入数据库
	ticket := &model.SsoOAuthTicket{
		Code:        str,                              // 数据库中存放未加密的code，客户端来认证的时候使用的是加密后的code，这样在验证code的时候将前端加密的进行解密判断是否与数据库中的相等即可
		RedirectURI: site.CallbackUrl,                 // 回调地址
		UserID:      userId,                           // 用户ID
		ExpiresAt:   time.Now().Add(10 * time.Second), // 票据的有效期为10秒
	}
	if err = dao.SSO.CreateAuthorizeCode(ticket); err != nil {
		return "", err
	}

	// 返回授权码
	redirectURI := fmt.Sprintf("%s?code=%s&state=%s", site.CallbackUrl, code, data.State)
	return redirectURI, nil
}

// GetToken 客户端Token获取
func (s *sso) GetToken(param *Token) (token *ResponseToken, err error) {

	var user *dao.UserInfoWithMenu

	// 客户端验证
	site, err := dao.Site.GetSite(param.ClientId)
	if err != nil {
		return nil, errors.New("client_id string is invalid")
	}
	if site.ClientSecret != param.ClientSecret {
		return nil, errors.New("client_secret string is invalid")
	}

	// 取票据验证（如果有数据则表明：1、Code存在，2、在有效期内，3、未使用）
	str, _ := utils.Decrypt(param.Code)
	ticket, err := dao.SSO.GetAuthorizeCode(str)
	if err != nil {
		return nil, errors.New("code string is invalid")
	}

	// 生成access_token（使用JWT生成，方便后续对Token进行校验）
	user, err = dao.User.GetUser(ticket.UserID)
	accessToken, err := middleware.GenerateJWT(uint(user.ID), user.Name, user.Username)
	if err != nil {
		return nil, err
	}

	token = &ResponseToken{
		AccessToken: accessToken,
		TokenType:   "bearer", // 固定值
		ExpiresIn:   3600,     // Token过期时间，这里和配置文件中的JWT过期时间保持一致，也可以独立配置
		Scope:       "openid", // 固定值
	}

	return token, err
}

// GetUserinfo 客户端获取用户信息
func (s *sso) GetUserinfo(token string) (user *ResponseUserinfo, err error) {

	// 验证Token
	mc, err := middleware.ValidateJWT(token)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	userinfo, err := dao.User.GetUser(mc.ID)

	user = &ResponseUserinfo{
		Id:          uint(userinfo.ID),
		Name:        userinfo.Name,
		Username:    userinfo.Username,
		Email:       userinfo.Email,
		PhoneNumber: userinfo.PhoneNumber,
	}

	return user, err
}
