package db

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/wonderivan/logger"
	"ops-api/config"
	"ops-api/global"
)

func MinioInit() error {

	// 读取配置信息
	endpoint := config.Conf.OSS.Endpoint
	accessKey := config.Conf.OSS.AccessKey
	secretKey := config.Conf.OSS.SecretKey

	// 客户端初始化
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: config.Conf.OSS.SSL,
	})
	if err != nil {
		return err
	}

	global.MinioClient = client
	logger.Info("Minio客户端初始化成功.")

	return nil
}
