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

func (r *minioRepo) UploadFile(ctx context.Context, objectName string, filePath string) (string, error) {
	uploadFile, err := r.data.minio.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: "text/plain; charset=utf-8"})
	if err != nil {
		return "", err
	}
	return uploadFile.Location, nil
}
