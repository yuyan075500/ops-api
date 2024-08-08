package dao

import (
	"errors"
	"ops-api/global"
	"ops-api/model"
	"time"
)

var SSO sso

type sso struct{}

// CreateAuthorizeCode 创建授权码（OAuth2.0）
func (l *sso) CreateAuthorizeCode(data *model.SsoOAuthTicket) (err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// CreateAuthorizeTicket 创建授权票据（CAS3.0）
func (l *sso) CreateAuthorizeTicket(data *model.SsoCASTicket) (err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// GetAuthorizeCode 仅获取有效授权码（OAuth2.0）
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

// GetAuthorizeTicket 仅获取有效票据（CAS3.0）
func (l *sso) GetAuthorizeTicket(st string) (data *model.SsoCASTicket, err error) {
	var ticket *model.SsoCASTicket

	// 仅获取有效票据（1、Ticket存在，2、在有效期内，3、未使用）
	now := time.Now()
	if err := global.MySQLClient.Where("ticket = ? AND expires_at > ? AND consumed_at IS NULL", st, now).First(&ticket).Error; err != nil {
		return nil, err
	}

	// 票据使用过后，进行使用标记（确保票据只能使用一次）
	if err := global.MySQLClient.Model(&ticket).Update("consumed_at", now).Error; err != nil {
		return nil, err
	}

	return ticket, nil
}
