package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// 全局变量
var (
	MinioClient  *minio.Client
	RedisClient  *redis.Client
	MySQLClient  *gorm.DB
	CasBinServer *casbin.Enforcer
)
