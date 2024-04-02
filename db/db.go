package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wonderivan/logger"
	"ops-api/config"
)

var (
	isInit bool
	GORM   *gorm.DB
	err    error
)

func Init() {
	//判断否已经初始化
	if isInit {
		return
	}

	//组装数据库连接配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	//建议数据库连接，并生成*gorm.DB对象
	GORM, err = gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println("数据库连接失败：" + err.Error())
	}

	//连接池相关配置
	GORM.DB().SetMaxIdleConns(config.MaxIdleConns)
	GORM.DB().SetMaxOpenConns(config.MaxOpenConns)
	GORM.DB().SetConnMaxLifetime(config.MaxLifeTime)

	isInit = true
	logger.Info("数据库初始化成功." + "\n")
}
