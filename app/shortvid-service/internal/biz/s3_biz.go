package biz

import (
	"context"
	"shortvid-backend/app/shortvid-service/internal/conf"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/go-kratos/kratos/v3/log"
)

type S3Repo interface {
	GetUploadSession(ctx context.Context) (*sts.AssumeRoleOutput, error)
	ListBuckets(ctx context.Context) ([]string, error)
}

type S3Usecase struct {
	conf *conf.S3
	repo S3Repo
}

func NewS3Usecase(conf *conf.S3, repo S3Repo) *S3Usecase {
	return &S3Usecase{conf: conf, repo: repo}
}

type UploadSession struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Endpoint        string
}

func (s *S3Usecase) GetUploadSession(ctx context.Context) (*UploadSession, error) {
	output, err := s.repo.GetUploadSession(ctx)
	if err != nil {
		log.ErrorContext(ctx, "get upload session failed", "error", err)
		return nil, err
	}
	return &UploadSession{
		AccessKeyID:     *output.Credentials.AccessKeyId,
		SecretAccessKey: *output.Credentials.SecretAccessKey,
		SessionToken:    *output.Credentials.SessionToken,
		Endpoint:        s.conf.Endpoint,
	}, nil
}

func (s *S3Usecase) ListBuckets(ctx context.Context) ([]string, error) {
	buckets, err := s.repo.ListBuckets(ctx)
	if err != nil {
		log.ErrorContext(ctx, "list buckets failed", "error", err)
		return nil, err
	}
	return buckets, nil
}
