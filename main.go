package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ops-api/config"
	"ops-api/controller"
	"ops-api/db"
	"ops-api/middleware"
)

func main() {

	// 配置初始化
	config.Init()

	// 初始化MySQL
	db.MySQLInit()

	// 初始Redis
	db.RedisInit()

	// 初始化Minio
	db.MinioInit()

	r := gin.Default()

	// 初始化跨域中间件
	r.Use(middleware.Cors())
	// 初始化登录中间件，其中IgnorePaths()方法可以忽略某些路由，支持前缀匹配
	r.Use(middleware.LoginBuilder().
		IgnorePaths("/login").
		IgnorePaths("/health").
		IgnorePaths("/swagger/").
		Build())

	// 注册路由
	controller.Router.InitRouter(r)

	// 启动服务
	err := r.Run(fmt.Sprintf("%v", config.Conf.Server))
	if err != nil {
		fmt.Println(err.Error())
	}
}
