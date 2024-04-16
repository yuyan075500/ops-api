package db

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/wonderivan/logger"
	"ops-api/config"
	"ops-api/global"
)

var MinioClient *minio.Client

func MinioInit() {
	endpoint := config.Conf.OSS.Endpoint
	accessKey := config.Conf.OSS.AccessKey
	secretKey := config.Conf.OSS.SecretKey
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: config.Conf.OSS.SSL,
	})
	if err != nil {
		logger.Error("Minio客户端初始化失败." + err.Error())
		return
	}

	global.MinioClient = minioClient
	logger.Info("Minio客户端初始化成功.")
}
