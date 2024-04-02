package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ops-api/config"
	"ops-api/controller"
	"ops-api/db"
)

func main() {
	r := gin.Default()

	// 初始化数据库
	db.Init()

	// 注册路由
	controller.Router.InitApiRouter(r)

	// 启动服务
	err := r.Run(config.ListenAddr)
	if err != nil {
		fmt.Println(err.Error())
	}
}
