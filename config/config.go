package config

import (
	"github.com/wonderivan/logger"
	"gopkg.in/yaml.v3"
	"os"
)

// Conf 全局变量
var Conf *Config

type Config struct {
	Server   string `yaml:"server"`
	Database MySQL  `yaml:"mysql"`
	JWT      JWT    `yaml:"jwt"`
	Redis    Redis  `yaml:"redis"`
	OSS      OSS    `yaml:"oss"`
}

type MySQL struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	DB           string `yaml:"db"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	MaxLifeTime  int    `yaml:"maxLifeTime"`
}

type Redis struct {
	Host     string `yaml:"host"`
	DB       int    `yaml:"db"`
	Password string `yaml:"password"`
}

type OSS struct {
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"accessKey"`
	SecretKey  string `yaml:"secretKey"`
	BucketName string `yaml:"bucketName"`
	SSL        bool   `yaml:"ssl"`
}

type JWT struct {
	Secret  string `yaml:"secret"`
	Expires int    `yaml:"expires"`
}

func Init() {
	// 加载配置文件
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		logger.Error("读取配置文件失败：%v", err)
		return
	}

	// 解析配置文件
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		logger.Error("配置文件解析失败: %v", err)
		return
	}

	// 将解析出来的配置赋值给全局变量
	Conf = &cfg
}
