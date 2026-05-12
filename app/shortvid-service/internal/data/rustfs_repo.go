package data

import (
	"context"
)

type rustFsRepo struct {
	data *Data
}

const (
	bucketName = "shortvid"
)

func NewRustFsRepo(data *Data) *rustFsRepo {
	return &rustFsRepo{data: data}
}

// UploadFile 上传文件到RustFs
func (r *rustFsRepo) UploadFile(ctx context.Context, objectName string, filePath string) error {
	return nil
}
