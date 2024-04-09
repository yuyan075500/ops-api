package db

import (
	"fmt"
	"github.com/wonderivan/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ops-api/config"
	"ops-api/model"
)

var (
	isInit bool
	GORM   *gorm.DB
	err    error
)

func Init() {
	// 判断否已经初始化
	if isInit {
		return
	}

	// 组装数据库连接配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Conf.Database.User,
		config.Conf.Database.Password,
		config.Conf.Database.Host,
		config.Conf.Database.Port,
		config.Conf.Database.DB,
	)

	// 建议数据库连接，并生成*gorm.DB对象
	GORM, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败：" + err.Error())
	}

	// 表迁移
	_ = GORM.SetupJoinTable(&model.AuthUser{}, "Groups", &model.AuthUserGroups{})
	_ = GORM.SetupJoinTable(&model.AuthGroup{}, "Permissions", &model.AuthGroupPermissions{})
	_ = GORM.AutoMigrate(
		&model.AuthUser{},
		&model.AuthGroup{},
		&model.AuthPermission{},
	)

	// 数据库连接池设置
	//DB, _ := GORM.DB()
	//DB.SetMaxIdleConns(config.Conf.Database.MaxIdleConns)
	//DB.SetMaxOpenConns(config.Conf.Database.MaxOpenConns)
	//DB.SetConnMaxLifetime(time.Duration(config.Conf.Database.MaxLifeTime))

	isInit = true
	logger.Info("数据库初始化成功." + "\n")
}
