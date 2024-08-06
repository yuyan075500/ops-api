package dao

import (
	"errors"
	"ops-api/global"
	"ops-api/model"
	"time"
)

var SSO sso

type sso struct{}

// CreateAuthorizeCode 创建授权码
func (l *sso) CreateAuthorizeCode(data *model.SsoOAuthTicket) (err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// GetAuthorizeCode 仅获取有效授权码
func (l *sso) GetAuthorizeCode(code string) (data *model.SsoOAuthTicket, err error) {
	var ticket *model.SsoOAuthTicket

	// 仅获取有效授权码（1、Code存在，2、在有效期内，3、未使用）
	now := time.Now()
	if err := global.MySQLClient.Where("code = ? AND expires_at > ? AND consumed_at IS NULL", code, now).First(&ticket).Error; err != nil {
		return nil, err
	}

	// 票据使用过后，进行使用标记（确保票据只能使用一次）
	if err := global.MySQLClient.Model(&ticket).Update("consumed_at", now).Error; err != nil {
		return nil, err
	}

	return ticket, nil
}
