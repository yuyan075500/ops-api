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
func MySQLInit() error {

	var (
		user      model.AuthUser
		siteGroup model.SiteGroup
		site1     model.Site
		site2     model.Site
		site3     model.Site
	)

	// 组装数据库连接配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
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
		return err
	}

	// 表迁移
	_ = client.AutoMigrate(
		&model.AuthUser{},
		&model.AuthGroup{},
		&model.SiteGroup{},
		&model.Tag{},
		&model.Site{},
		&model.Menu{},
		&model.SubMenu{},
		&model.SystemPath{},
		&model.LogSMS{},
		&model.LogLogin{},
		&model.SsoOAuthTicket{},
		&model.SsoCASTicket{},
		&model.ScheduledTask{},
		&model.ScheduledTaskExecLog{},
		&model.Account{},
	)

	// 设置数据库连接池
	DB, _ := client.DB()
	DB.SetMaxIdleConns(config.Conf.MySQL.MaxIdleConns)
	DB.SetMaxOpenConns(config.Conf.MySQL.MaxOpenConns)
	DB.SetConnMaxLifetime(time.Duration(config.Conf.MySQL.MaxLifeTime) * time.Second)

	global.MySQLClient = client
	logger.Info("MySQL数据库初始化成功.")

	// 创建超级用户
	result := global.MySQLClient.FirstOrCreate(&user, model.AuthUser{
		Name:     "管理员",
		Username: "admin",
		IsActive: true,
		Password: "admin@123...",
		WwId:     nil,
	})

	if result.RowsAffected == 0 {
		logger.Warn("admin用户已存在.")
	} else {
		logger.Info("admin用户创建成功.")
		// 创建初始站点
		global.MySQLClient.FirstOrCreate(&siteGroup, model.SiteGroup{
			Name: "信息化公用",
		})
		global.MySQLClient.FirstOrCreate(&site1, model.Site{
			Name:        "密码重置",
			Description: "统一认证平台密码自助更改平台，支持本地、Windows AD和OpenLDAP用户密码修改",
			Address:     fmt.Sprintf("%s/reset_password", config.Conf.ExternalUrl),
			SSO:         false,
			SiteGroupID: siteGroup.ID,
		})
		global.MySQLClient.FirstOrCreate(&site2, model.Site{
			Name:        "接口文档",
			Description: "统一认证平台Swagger接口文档，生产环境建议关闭",
			Address:     fmt.Sprintf("%s/swagger/index.html", config.Conf.ExternalUrl),
			SSO:         false,
			SiteGroupID: siteGroup.ID,
		})
		global.MySQLClient.FirstOrCreate(&site3, model.Site{
			Name:        "站点导航",
			Description: "统一认证平台站点导航，无需认证，可在后台进行编辑",
			Address:     fmt.Sprintf("%s/sites", config.Conf.ExternalUrl),
			SSO:         false,
			SiteGroupID: siteGroup.ID,
		})
	}

	return nil
}
