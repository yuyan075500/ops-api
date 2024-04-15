package db

import (
	"github.com/go-redis/redis"
	"github.com/wonderivan/logger"
	"ops-api/config"
)

var rdb *redis.Client

// RedisInit Redis初始化
func RedisInit() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Host,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB,
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		logger.Info("Redis客户端初始化失败." + err.Error())
		return
	}

	logger.Info("Redis客户端初始化成功.")
}
