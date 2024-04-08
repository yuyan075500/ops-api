package db

import (
	"fmt"
	"github.com/wonderivan/logger"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ops-api/model"
)

var (
	isInit bool
	GORM   *gorm.DB
	err    error
)

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	DB           string `yaml:"db"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	MaxLifeTime  int    `yaml:"maxLifeTime"`
}

type Config struct {
	Database DatabaseConfig `yaml:"mysql"`
}

func Init(config []byte) {
	// 判断否已经初始化
	if isInit {
		return
	}

	// 读取配置
	var conf Config
	_ = yaml.Unmarshal(config, &conf)

	// 组装数据库连接配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DB,
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
	//DB.SetMaxIdleConns(conf.Database.MaxIdleConns)
	//DB.SetMaxOpenConns(conf.Database.MaxOpenConns)
	//DB.SetConnMaxLifetime(time.Duration(conf.Database.MaxLifeTime))

	isInit = true
	logger.Info("数据库初始化成功." + "\n")
}
