package middleware

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/wonderivan/logger"
	"ops-api/config"
	"ops-api/service"
)

type LDAPServer struct {
	Conn   *ldap.Conn
	Config LDAPConfig
}

type LDAPConfig struct {
	Addr             string
	BindUserDN       string
	BindUserPassword string
	SearchDN         string
}

// CreateLDAPService 建立LDAP连接
func CreateLDAPService() (*LDAPServer, error) {
	conf := LDAPConfig{
		Addr:             config.Conf.LDAP.Host,
		BindUserDN:       config.Conf.LDAP.BindUserDN,
		BindUserPassword: config.Conf.LDAP.BindUserPassword,
		SearchDN:         config.Conf.LDAP.SearchDN,
	}

	conn, err := ldap.DialURL(conf.Addr, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return nil, err
	}
	_, err = conn.SimpleBind(&ldap.SimpleBindRequest{
		Username: conf.BindUserDN,
		Password: conf.BindUserPassword,
	})
	if err != nil {
		return nil, err
	}

	return &LDAPServer{Conn: conn, Config: conf}, nil
}

// LDAPUserAuthentication 用户认证
func (l *LDAPServer) LDAPUserAuthentication(username, password string) (data *service.UserCreate, err error) {

	// 查找用户
	searchRequest := ldap.NewSearchRequest(
		config.Conf.LDAP.SearchDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(sAMAccountName=%s)", username),
		[]string{},
		nil,
	)
	searchResult, err := l.Conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	if len(searchResult.Entries) != 1 {
		return nil, errors.New("用户不存在")
	}

	// 验证用户密码是否正确
	userDN := searchResult.Entries[0].DN
	err = l.Conn.Bind(userDN, password)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		return nil, errors.New("用户密码错误或账号异常")
	}

	// 返回用户信息
	return &service.UserCreate{
		Name:        searchResult.Entries[0].GetAttributeValue("cn"),
		Username:    searchResult.Entries[0].GetAttributeValue("sAMAccountName"),
		Password:    password,
		PhoneNumber: searchResult.Entries[0].GetAttributeValue("telephoneNumber"),
		Email:       searchResult.Entries[0].GetAttributeValue("mail"),
		UserFrom:    "AD域",
	}, nil
}
