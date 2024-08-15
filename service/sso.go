package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/LoginRadius/go-saml"
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

// SAMLRequest SAML2客户端授权请求参数
type SAMLRequest struct {
	SAMLRequest string `form:"SAMLRequest" binding:"required"` // SAMLRequest数据，通常该数据是DEFLATE压缩 + base64编码，获取此数据需要进行DEFLATE解压缩 + base64解码
	RelayState  string `form:"RelayState"`                     // SP的状态信息，防止跨站请求伪造攻击，功能与OAuth2.0客户端的state功能相同
	SigAlg      string `form:"SigAlg"`                         // 签名使用的算法
	Signature   string `form:"Signature"`                      // 签名，用于验证SP的身份，但需要配置SP的公钥
}

// ParseSPMetadata 获取SP Metadata信息请求参数
type ParseSPMetadata struct {
	SPMetadataURL string `json:"sp_metadata_url" binding:"required"`
}

// SPMetadata 返回给前端的SP Metadata数据
type SPMetadata struct {
	EntityID    string `json:"entity_id"`
	Certificate string `json:"certificate"`
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

// SAMLRequestData SAMLRequest数据绑定结构体
type SAMLRequestData struct {
	XMLName                     xml.Name     `xml:"urn:oasis:names:tc:SAML:2.0:protocol AuthnRequest"`
	AssertionConsumerServiceURL string       `xml:"AssertionConsumerServiceURL,attr"`
	Destination                 string       `xml:"Destination,attr"`
	ID                          string       `xml:"ID,attr"`
	IssueInstant                string       `xml:"IssueInstant,attr"`
	ProtocolBinding             string       `xml:"ProtocolBinding,attr"`
	Version                     string       `xml:"Version,attr"`
	Issuer                      Issuer       `xml:"urn:oasis:names:tc:SAML:2.0:assertion Issuer"`
	NameIDPolicy                NameIDPolicy `xml:"NameIDPolicy"`
}
type Issuer struct {
	Value string `xml:",chardata"`
}
type NameIDPolicy struct {
	AllowCreate     string `xml:"AllowCreate,attr"`
	Format          string `xml:"Format,attr"`
	SPNameQualifier string `xml:"SPNameQualifier,attr"`
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

// GetIdPMetadata 获取SAML2 IDP Metadata
func (s *sso) GetIdPMetadata() (metadata string, err error) {

	// 获取证书
	cert, err := utils.LoadIdpCertificate()
	if err != nil {
		return "", err
	}

	// 创建IDP实例
	idp := saml.IdentityProvider{
		IsIdpInitiated:       false,                      // 是否是IdP Initiated模式，true：表示认证请求是通过IdP发起的，false：表示认证请求是客户端（SP）发起的
		Issuer:               "https://ops-test.50yc.cn", // IDP实体，默认为当前服务器地址
		IDPCert:              base64.StdEncoding.EncodeToString(cert.Raw),
		NameIdentifierFormat: saml.AttributeFormatUnspecified,
	}

	// 添加单点登录接口信息
	idp.AddSingleSignOnService(saml.MetadataBinding{
		Binding:  saml.HTTPRedirectBinding,         // 由于IDP是前后端分离架构，所以这里使用HTTPRedirectBinding
		Location: "https://ops-test.50yc.cn/login", // 单点登录接口地址
	})

	// 添加单点登出接口信息（不支持：如果支持单点登出则可以添加此信息到元数据中）
	//idp.AddSingleSignOutService(saml.MetadataBinding{
	//	Binding:  saml.HTTPPostBinding,
	//	Location: "https://ops-test.50yc.cn/logout",
	//}

	// 添加IDP组件相关信息
	idp.AddOrganization(saml.Organization{
		OrganizationDisplayName: "运维平台", // 组织显示名称
		OrganizationName:        "OPS",  // 组织正式名称
		OrganizationURL:         "https://ops-test.50yc.cn",
	})

	// 添加主要联系人信息
	idp.AddContactPerson(saml.ContactPerson{
		ContactType:  "technical", // 联系人类型包含：technical（技术联系人）、support（支持联系人）、administrative（行政联系人）、billing（财务联系人）、other（其它）
		EmailAddress: "zhangs@ops.cn",
		GivenName:    "三",
		SurName:      "张",
	})

	// 添加其它联系人信息
	//persons := []saml.ContactPerson{
	//	{
	//		ContactType:  "support",
	//		EmailAddress: "support@ops.cn",
	//		GivenName:    "四",
	//		SurName:      "李",
	//	},
	//}
	//idp.AddContactPersons(persons...)

	// 生成metadata元数据
	metadata, msg := idp.MetaDataResponse()
	if msg != nil {
		return "", msg.Error
	}

	return metadata, nil
}

// ParseSPMetadata SP Metadata解析
func (s *sso) ParseSPMetadata(metadataUrl string) (data *SPMetadata, err error) {

	metadata, err := utils.ParseSPMetadata(metadataUrl)
	if err != nil {
		return nil, err
	}

	// 提取IDP的签名证书
	var signingCertData string
	for _, keyDescriptor := range metadata.SPSSODescriptor.KeyDescriptors {
		if keyDescriptor.Use == "signing" {
			signingCertData = keyDescriptor.KeyInfo.X509Data.X509Certificate
			break
		}
	}
	if signingCertData == "" {
		return nil, errors.New("未找到签名证书")
	}

	return &SPMetadata{
		Certificate: signingCertData,
		EntityID:    metadata.EntityID,
	}, nil
}

// SPAuthorize SP授权
func (s *sso) SPAuthorize(samlRequest *SAMLRequest, userId uint) (err error) {

	var data *SAMLRequestData

	// 获取SAMLRequest数据
	samlRequestRaw, err := utils.ParseSAMLRequest(samlRequest.SAMLRequest)
	if err != nil {
		return err
	}

	// 将SAMLRequest数据到结构体
	if err := xml.Unmarshal([]byte(samlRequestRaw.String()), &data); err != nil {
		return err
	}

	// 获取SP应用
	site, err := dao.Site.GetSamlSite(data.Issuer.Value)
	if err != nil {
		return err
	}

	// 判断用户是否有权限访问
	if !site.AllOpen {
		if !dao.Site.IsUserInSite(userId, site) {
			return errors.New("您无权访问该应用")
		}
	}

	// 签名验证
	//if err := utils.VerifySignature(samlRequestRaw, site.Certificate, samlRequest.Signature, samlRequest.SigAlg); err != nil {
	//	return err
	//}

	// 生成SAMLResponse

	return nil
}
