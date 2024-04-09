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

	// 初始化数据库
	db.Init()

	r := gin.Default()

	// 初始化中间件
	r.Use(middleware.LoginBuilder().
		IgnorePaths("/login").
		Build())

	// 注册路由
	controller.Router.InitApiRouter(r)

	// 启动服务
	err := r.Run(fmt.Sprintf("%v", config.Conf.Server))
	if err != nil {
		fmt.Println(err.Error())
	}
}
