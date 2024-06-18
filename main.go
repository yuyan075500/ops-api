package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"ops-api/config"
	"ops-api/controller"
	"ops-api/db"
	"ops-api/middleware"
)

func main() {

	// 配置初始化
	config.Init()

	// 初始化MySQL
	if err := db.MySQLInit(); err != nil {
		logger.Error("ERROR：", err.Error())
		return
	}

	// 初始Redis
	if err := db.RedisInit(); err != nil {
		logger.Error("ERROR：", err.Error())
		return
	}

	// 初始化Minio
	if err := db.MinioInit(); err != nil {
		logger.Error("ERROR：", err.Error())
		return
	}

	// 初始化CasBin权限
	if err := middleware.CasBinInit(); err != nil {
		logger.Error("ERROR：", err.Error())
		return
	}

	r := gin.Default()

	// 加载跨域中间件
	r.Use(middleware.Cors())
	// 加载登录中间件，其中IgnorePaths()方法可以忽略某些路由，支持前缀匹配
	r.Use(middleware.LoginBuilder().
		IgnorePaths("/login").
		IgnorePaths("/health").
		IgnorePaths("/swagger/").
		IgnorePaths("/api/v1/sms/callback").
		Build())
	// 加载权限中间件
	r.Use(middleware.PermissionCheck())

	// 注册路由
	controller.Router.InitRouter(r)

	// 启动服务
	if err := r.Run(fmt.Sprintf("%v", config.Conf.Server)); err != nil {
		logger.Error("ERROR：", err.Error())
	}
}
