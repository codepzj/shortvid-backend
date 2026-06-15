package storage

import (
	"context"
	"log"

	"shortvid-backend/app/shortvid-service/internal/conf"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type S3Data struct {
	S3Conf          *conf.S3
	Client          *s3.Client        // s3对象存储客户端
	PresignedClient *s3.PresignClient // s3临时凭证客户端
	StsClient       *sts.Client       // sts客户端
}

// NewS3 使用aws的go sdk连接S3对象存储
func NewS3(conf *conf.S3) *S3Data {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(conf.GetRegion()),
		config.WithBaseEndpoint(conf.GetEndpoint()),
	)
	if err != nil {
		log.Fatalf("aws sdk load config failed: %v", err)
	}

	// 注意
	// s3客户端应该使用控制台创建的ak和sk
	// sts客户端实际上使用登录s3 web对应的账号和密码

	// s3客户端
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = conf.GetUsePathStyle() // minio,rustfs等需要使用路径风格
		o.Credentials = credentials.NewStaticCredentialsProvider(conf.GetS3AccessKey(), conf.GetS3SecretKey(), "")
	})
	presignedClient := s3.NewPresignClient(client)

	// sts客户端
	stsClient := sts.NewFromConfig(cfg, func(o *sts.Options) {
		o.Credentials = credentials.NewStaticCredentialsProvider(conf.GetStsAccessKey(), conf.GetStsSecretKey(), "")
	})

	log.Println("S3 data client initialized successfully...")

	return &S3Data{
		S3Conf:          conf,
		Client:          client,
		PresignedClient: presignedClient,
		StsClient:       stsClient,
	}
}
