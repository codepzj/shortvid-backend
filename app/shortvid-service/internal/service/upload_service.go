package service

import (
	"shortvid-backend/app/shortvid-service/internal/data"

)

type UploadService struct {
	data *data.Data
}

func NewUploadService(data *data.Data) *UploadService {
	return &UploadService{data: data}
}

