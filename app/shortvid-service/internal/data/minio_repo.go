package data

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type minioRepo struct {
	data *Data
}

const (
	bucketName = "shortvid"
)

func NewMinioRepo(data *Data) *minioRepo {
	return &minioRepo{data: data}
}

// UploadFile 上传文件到Minio
func (r *minioRepo) UploadFile(ctx context.Context, objectName string, filePath string) error {
	_, err := r.data.minio.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
	return err
}
