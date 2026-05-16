package biz

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/go-kratos/kratos/v2/log"
)

type S3Repo interface {
	GetUploadPresignedURL(ctx context.Context) (string, error)
	GetUploadSession(ctx context.Context) (*sts.AssumeRoleOutput, error)
	ListBuckets(ctx context.Context) ([]string, error)
}

type S3Usecase struct {
	logger *log.Helper
	repo   S3Repo
}

func NewS3Usecase(logger log.Logger, repo S3Repo) *S3Usecase {
	return &S3Usecase{logger: log.NewHelper(logger), repo: repo}
}

func (s *S3Usecase) GetUploadPresignedURL(ctx context.Context) (string, error) {
	return s.repo.GetUploadPresignedURL(ctx)
}

type UploadSession struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

func (s *S3Usecase) GetUploadSession(ctx context.Context) (*UploadSession, error) {
	output, err := s.repo.GetUploadSession(ctx)
	if err != nil {
		s.logger.Error("get upload session failed", "error", err)
		return nil, err
	}
	return &UploadSession{
		AccessKeyID:     *output.Credentials.AccessKeyId,
		SecretAccessKey: *output.Credentials.SecretAccessKey,
		SessionToken:    *output.Credentials.SessionToken,
	}, nil
}

func (s *S3Usecase) ListBuckets(ctx context.Context) ([]string, error) {
	buckets, err := s.repo.ListBuckets(ctx)
	if err != nil {
		s.logger.Error("list buckets failed", "error", err)
		return nil, err
	}
	return buckets, nil
}
