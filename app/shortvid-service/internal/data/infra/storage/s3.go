package storage

import (
	"context"
	"log"

	"shortvid-backend/app/shortvid-service/internal/conf"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type S3Data struct {
	Client    *s3.Client  // s3对象存储客户端
	StsClient *sts.Client // sts临时凭证客户端
}

// NewS3 使用aws的go sdk连接RustFS对象存储
func NewS3(conf *conf.RustFs) *S3Data {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.GetAccessKey(), conf.GetSecretKey(), "")),
	)
	if err != nil {
		log.Fatalf("aws sdk load config failed: %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(conf.GetEndpoint())
		o.UsePathStyle = true
	})
	stsClient := sts.NewFromConfig(cfg, func(o *sts.Options) {
		o.BaseEndpoint = aws.String(conf.GetEndpoint())
	})

	log.Println("S3 client and stsClient initialized successfully...")

	return &S3Data{
		Client:    client,
		StsClient: stsClient,
	}
}
