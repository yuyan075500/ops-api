package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"gopkg.in/yaml.v3"
	"ops-api/controller"
	"ops-api/db"
	"ops-api/middleware"
	"ops-api/utils"
)

type ServerConfig struct {
	ListenAddress string `yaml:"server"`
}

func main() {
	r := gin.Default()

	// 加载配置文件
	config, err := utils.ReadFile("config/conf.yaml")
	if err != nil {
		logger.Error("加载配置文件失败.")
	}

	// 初始化数据库
	db.Init(config)

	// 初始化中间件
	r.Use(middleware.LoginBuilder().
		IgnorePaths("/user/login").
		Build())

	// 注册路由
	controller.Router.InitApiRouter(r)

	// 读取配置
	var server ServerConfig
	_ = yaml.Unmarshal(config, &server)

	// 启动服务
	err = r.Run(fmt.Sprintf("%v", server.ListenAddress))
	if err != nil {
		fmt.Println(err.Error())
	}
}
