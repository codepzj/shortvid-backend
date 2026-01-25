package storage

import (
	"log"
	"shortvid-backend/app/shortvid-service/internal/conf"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient(conf *conf.Minio) *minio.Client {
	minioClient, err := minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		log.Fatalf("Connect Minio failed: %v", err)
	}
	log.Printf("Minio connect success...")
	return minioClient
}
