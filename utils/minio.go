package utils

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
	"ops-api/config"
	"ops-api/global"
)

func FileUpload(fileName, ContentType string, file io.Reader, fileSize int64) (err error) {
	_, err = global.MinioClient.PutObject(context.Background(), config.Conf.OSS.BucketName, fileName, file, fileSize, minio.PutObjectOptions{
		ContentType: ContentType,
	})
	if err != nil {
		return err
	}

	return nil
}
