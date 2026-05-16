package biz

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type S3Repo interface {
	GetUploadSession(ctx context.Context) (*sts.GetSessionTokenOutput, error)
}

type S3Usecase struct {
	repo S3Repo
}

func NewS3Usecase(repo S3Repo) *S3Usecase {
	return &S3Usecase{repo: repo}
}

type S3UploadSession struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

func (s *S3Usecase) GetUploadSession(ctx context.Context) (*S3UploadSession, error) {
	session, err := s.repo.GetUploadSession(ctx)
	if err != nil {
		return nil, err
	}

	if session.Credentials == nil {
		return nil, errors.New("session credentials is nil")
	}
	return &S3UploadSession{
		AccessKeyId:     *session.Credentials.AccessKeyId,
		SecretAccessKey: *session.Credentials.SecretAccessKey,
		SessionToken:    *session.Credentials.SessionToken,
	}, nil
}
