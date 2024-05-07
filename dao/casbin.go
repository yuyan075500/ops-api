package dao

import (
	"errors"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
	"ops-api/model"
)

var CasBin casbin

type casbin struct{}

func (u *casbin) AddRole(tx *gorm.DB, data *model.CasbinRule) (err error) {
	if err := tx.Create(&data).Error; err != nil {
		logger.Error("新增失败：", err)
		return errors.New("新增失败：" + err.Error())
	}
	return nil
}
