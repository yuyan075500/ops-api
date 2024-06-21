package service

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/encoding/unicode"
	"ops-api/config"
	"strings"
)

var AD ad

type ad struct{}

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

// Connect 建立LDAP连接
func (a *ad) Connect() (*LDAPServer, error) {
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

// LDAPUserSearch 根据用户名查找用户信息
func (a *ad) LDAPUserSearch(username string) (result *ldap.SearchResult, err error) {

	// 建立LDAP连接
	l, err := a.Connect()
	if err != nil {
		return nil, err
	}

	// 查找用户
	searchDN := strings.Split(config.Conf.LDAP.SearchDN, "&")
	for _, dn := range searchDN {

		// 构建查找请求
		searchRequest := ldap.NewSearchRequest(
			dn,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0,
			0,
			false,
			fmt.Sprintf("(&(objectClass=person)(sAMAccountName=%s))", username),
			[]string{},
			nil,
		)

		// 执行查找
		searchResult, err := l.Conn.Search(searchRequest)
		if err != nil {
			return nil, err
		}

		// 如果没有找到用户，则继续在下一个DN查找
		if len(searchResult.Entries) == 0 {
			continue
		}

		// 返回用户信息
		return searchResult, nil
	}

	return nil, errors.New("用户不存在")
}

// LDAPUserAuthentication 用户认证
func (a *ad) LDAPUserAuthentication(username, password string) (result *ldap.SearchResult, err error) {

	// 建立LDAP连接
	l, err := a.Connect()
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	searchResult, err := a.LDAPUserSearch(username)
	if err != nil {
		return nil, err
	}

	// 密码认证
	userDN := searchResult.Entries[0].DN
	err = l.Conn.Bind(userDN, password)
	if err != nil {
		return nil, errors.New("用户密码错误或账号异常")
	}

	// 返回用户信息
	return searchResult, nil
}

// LDAPUserResetPassword 重置用户密码
func (a *ad) LDAPUserResetPassword(username, password string) (err error) {
	// 建立LDAP连接
	l, err := a.Connect()
	if err != nil {
		return err
	}

	// 对密码进行utf16编码
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	pwdEncoded, _ := utf16.NewEncoder().String("\"" + password + "\"")

	// 获取用户信息
	searchResult, err := a.LDAPUserSearch(username)
	if err != nil {
		return err
	}

	// 构建修改请求
	userDN := searchResult.Entries[0].DN
	req := ldap.NewModifyRequest(userDN, []ldap.Control{})

	// 修改密码
	req.Replace("unicodePwd", []string{pwdEncoded})

	// 修改用户账户状态
	//req.Replace("userAccountControl", []string{fmt.Sprintf("%d", 512)})

	// 执行修改请求
	// 注意：修改用户密码需要确保BindUserDN账号具备修改用户密码权限，以及需要使用ldaps方式连接，ldaps默认端口号为636，如：ldaps://192.168.200.13:636
	if err := l.Conn.Modify(req); err != nil {
		return err
	}
	return nil
}

// LDAPUserSync 用户同步
func (a *ad) LDAPUserSync() (users []UserSync, err error) {

	var userList []UserSync

	// 建立LDAP连接
	l, err := a.Connect()
	if err != nil {
		return nil, err
	}

	// 查找所有用户
	searchDN := strings.Split(config.Conf.LDAP.SearchDN, "&")
	for _, dn := range searchDN {
		// 构建查找请求
		searchRequest := ldap.NewSearchRequest(
			dn, // 指定查找范围
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0,
			0,
			false,
			fmt.Sprintf("(objectClass=person)"), // 指定过滤条件：类型为用户
			[]string{},
			nil,
		)
		// 执行查找
		searchResult, err := l.Conn.Search(searchRequest)
		if err != nil {
			return nil, err
		}
		// 获取查结果
		for _, value := range searchResult.Entries {

			// 判断用户状态
			// userAccountControl的值说明：514 禁用，512 启用，66050 禁用+密码永不过期，66048 启用+密码永不过期
			var isActive bool
			if value.GetAttributeValue("userAccountControl") == "514" || value.GetAttributeValue("userAccountControl") == "66050" {
				isActive = false
			} else {
				isActive = true
			}

			// 获取用户信息
			userInfo := &UserSync{
				Name:        value.GetAttributeValue("cn"),
				Username:    value.GetAttributeValue("sAMAccountName"),
				Password:    "",
				IsActive:    isActive,
				PhoneNumber: value.GetAttributeValue("telephoneNumber"),
				Email:       value.GetAttributeValue("mail"),
				UserFrom:    "AD域",
			}
			// 将用户信息追加到结构体
			userList = append(userList, *userInfo)
		}
	}

	return userList, nil
}
