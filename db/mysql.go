package db

import (
	"fmt"
	"github.com/wonderivan/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"ops-api/config"
	"ops-api/global"
	"ops-api/model"
	"time"
)

// MySQLInit MySQL初始化
func MySQLInit() {

	// 组装数据库连接配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Conf.MySQL.User,
		config.Conf.MySQL.Password,
		config.Conf.MySQL.Host,
		config.Conf.MySQL.Port,
		config.Conf.MySQL.DB,
	)

	// 建立数据库连接，并生成*gorm.DB对象
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger2.Default.LogMode(logger2.Silent),
	})
	if err != nil {
		logger.Error("ERROR：", err.Error())
		return
	}

	// 表迁移
	_ = client.AutoMigrate(
		&model.AuthUser{},
		&model.AuthGroup{},
		&model.Menu{},
		&model.SubMenu{},
		&model.SystemPath{},
	)

	// 设置数据库连接池
	DB, _ := client.DB()
	DB.SetMaxIdleConns(config.Conf.MySQL.MaxIdleConns)
	DB.SetMaxOpenConns(config.Conf.MySQL.MaxOpenConns)
	DB.SetConnMaxLifetime(time.Duration(config.Conf.MySQL.MaxLifeTime) * time.Second)

	global.MySQLClient = client
	logger.Info("MySQL数据库初始化成功.")
}
