package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"ops-api/dao"
	"ops-api/middleware"
	"ops-api/model"
	"ops-api/utils"
	"strings"
	"time"
)

var SSO sso

type sso struct{}

// OAuthAuthorize OAuth2.0客户端获取授权请求参数
type OAuthAuthorize struct {
	ResponseType string `json:"response_type" binding:"required"`
	ClientId     string `json:"client_id" binding:"required"`
	RedirectURI  string `json:"redirect_uri"`
	State        string `json:"state"`
	Scope        string `json:"scope"`
}

// CASAuthorize CAS3.0客户端获取授权请求参数
type CASAuthorize struct {
	Service string `form:"service" binding:"required"`
}

// Token OAuth2.0客户端获取token请求参数
type Token struct {
	GrantType    string `form:"grant_type"`
	Code         string `form:"code"`
	ClientId     string `form:"client_id"`
	RedirectURI  string `form:"redirect_uri"`
	ClientSecret string `form:"client_secret"`
}

// CASServiceValidate CAS3.0客户端票据校验请求参数
type CASServiceValidate struct {
	Service string `form:"service" binding:"required"`
	Ticket  string `form:"ticket" binding:"required"`
}

// ResponseToken 返回给OAuth2.0客户端客户端的Token信息
type ResponseToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// CASServiceResponse CAS3.0客户端返回给客户端的用户信息
type CASServiceResponse struct {
	XMLName               xml.Name               `xml:"cas:serviceResponse"`
	Xmlns                 string                 `xml:"xmlns:cas,attr"`
	AuthenticationSuccess *AuthenticationSuccess `xml:"cas:authenticationSuccess"`
}
type AuthenticationSuccess struct {
	User       string     `xml:"cas:user"`
	Attributes Attributes `xml:"cas:attributes"`
}
type Attributes struct {
	Id          uint   `xml:"id"`
	Name        string `xml:"name"`
	Username    string `xml:"username"`
	Email       string `xml:"email"`
	PhoneNumber string `xml:"phone_number"`
}

// ResponseUserinfo 返回给客户端的用户信息
type ResponseUserinfo struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

// GetCASAuthorize CAS3.0客户端授权
func (s *sso) GetCASAuthorize(data *CASAuthorize, userId uint, username string) (callbackUrl string, err error) {

	// 获取客户端应用
	site, err := dao.Site.GetCASSite(data.Service)
	if err != nil {
		return "", err
	}

	// 判断用户是否有权限访问
	if !site.AllOpen {
		if !dao.Site.IsUserInSite(userId, site) {
			return "", errors.New("您无权访问该应用")
		}
	}

	// 生成票据（固定格式）
	st := fmt.Sprintf("ST-%d-%s", time.Now().Unix(), username)

	// 使用HMAC SHA-256对票据进行签名（使用的是JWT的Secret）
	mac := hmac.New(sha256.New, []byte("config.Conf.JWT.Secret"))
	mac.Write([]byte(st))
	signature := hex.EncodeToString(mac.Sum(nil))

	// 将授权票据写入数据库
	st = fmt.Sprintf("%s-%s", st, signature)
	ticket := &model.SsoCASTicket{
		Ticket:    st,                               // 票据信息
		Service:   site.CallbackUrl,                 // 回调地址
		UserID:    userId,                           // 用户ID
		ExpiresAt: time.Now().Add(10 * time.Second), // 票据的有效期为10秒
	}
	if err = dao.SSO.CreateAuthorizeTicket(ticket); err != nil {
		return "", err
	}

	// 返回票据
	redirectURI := fmt.Sprintf("%s?ticket=%s", site.CallbackUrl, st)
	return redirectURI, nil
}

// GetOAuthAuthorize OAuth2.0客户端授权
func (s *sso) GetOAuthAuthorize(data *OAuthAuthorize, userId uint) (callbackUrl string, err error) {

	// 获取客户端应用
	site, err := dao.Site.GetOAuthSite(data.ClientId)
	if err != nil {
		return "", err
	}

	// 判断用户是否有权限访问
	if !site.AllOpen {
		if !dao.Site.IsUserInSite(userId, site) {
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

// GetToken OAuth2.0客户端Token获取
func (s *sso) GetToken(param *Token) (token *ResponseToken, err error) {

	var user *dao.UserInfoWithMenu

	// 客户端验证
	site, err := dao.Site.GetOAuthSite(param.ClientId)
	if err != nil {
		return nil, errors.New("client_id string is invalid")
	}
	if site.ClientSecret != param.ClientSecret {
		return nil, errors.New("client_secret string is invalid")
	}

	// 获取Code（如果有数据则表明：1、Code存在，2、在有效期内，3、未使用）
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

// ServiceValidate CAS3.0客户端票据校验
func (s *sso) ServiceValidate(param *CASServiceValidate) (data *CASServiceResponse, err error) {
	// 客户端验证
	_, err = dao.Site.GetCASSite(param.Service)
	if err != nil {
		return nil, errors.New("service string is invalid")
	}

	// 获取票据（如果有数据则表明：1、Code存在，2、在有效期内，3、未使用）
	ticketInfo, err := dao.SSO.GetAuthorizeTicket(param.Ticket)
	if err != nil {
		return nil, errors.New("ticket string is invalid")
	}

	// 分离票据
	parts := strings.Split(param.Ticket, "-")

	// 票据验证：结构验证
	if len(parts) != 4 {
		return nil, errors.New("ticket string is invalid")
	}

	// 获取票据本体
	ticket := fmt.Sprintf("%s-%s-%s", parts[0], parts[1], parts[2])
	// 获取票据签名
	signature := parts[3]

	// 生成新的签名
	mac := hmac.New(sha256.New, []byte("config.Conf.JWT.Secret"))
	mac.Write([]byte(ticket))
	newSignature := hex.EncodeToString(mac.Sum(nil))

	// 票据验证：比较签名
	if !hmac.Equal([]byte(newSignature), []byte(signature)) {
		return nil, errors.New("ticket string is invalid")
	}

	// 获取用户信息
	user, err := dao.User.GetUser(ticketInfo.UserID)
	if err != nil {
		return nil, err
	}

	return &CASServiceResponse{
		Xmlns: "http://www.yale.edu/tp/cas",
		AuthenticationSuccess: &AuthenticationSuccess{
			User: user.Name,
			Attributes: Attributes{
				Id:          uint(user.ID),
				Email:       user.Email,
				Name:        user.Name,
				PhoneNumber: user.PhoneNumber,
				Username:    user.Username,
			},
		},
	}, nil
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
