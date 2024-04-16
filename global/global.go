package global

import (
	"github.com/go-redis/redis"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

var (
	MinioClient *minio.Client
	RedisClient *redis.Client
	MySQLClient *gorm.DB
)
