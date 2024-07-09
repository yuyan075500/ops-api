package utils

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
	"net/url"
	"ops-api/config"
	"ops-api/global"
	"time"
)

// FileUpload 上传对象
func FileUpload(fileName, ContentType string, file io.Reader, fileSize int64) (err error) {
	_, err = global.MinioClient.PutObject(context.Background(), config.Conf.OSS.BucketName, fileName, file, fileSize, minio.PutObjectOptions{
		ContentType: ContentType,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetPresignedURL 获取临时访问链接
func GetPresignedURL(fileName string, expiryTime time.Duration) (url *url.URL, err error) {

	presignedURL, err := global.MinioClient.PresignedGetObject(context.Background(), config.Conf.OSS.BucketName, fileName, expiryTime, nil)
	if err != nil {
		return nil, err
	}

	return presignedURL, nil
}

// StatObject 获取对象信息
func StatObject(objectName string) (objectInfo *minio.ObjectInfo, err error) {

	info, err := global.MinioClient.StatObject(context.Background(), config.Conf.OSS.BucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &info, nil
}
