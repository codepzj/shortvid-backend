package data

import (
	"context"
	"shortvid-backend/app/shortvid-service/internal/biz"
	"shortvid-backend/app/shortvid-service/internal/conf"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/go-kratos/kratos/v2/log"
)

type s3Repo struct {
	conf   *conf.S3
	data   *Data
	logger *log.Helper
}

func NewS3Repo(conf *conf.S3, data *Data, logger log.Logger) biz.S3Repo {
	return &s3Repo{
		conf:   conf,
		data:   data,
		logger: log.NewHelper(logger),
	}
}

func (r *s3Repo) GetUploadSession(ctx context.Context) (*sts.AssumeRoleOutput, error) {
	output, err := r.data.s3.StsClient.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         &r.conf.StsRoleArn,
		RoleSessionName: &r.conf.StsSessionName,
		Policy:          &r.conf.StsPolicy,
		DurationSeconds: &r.conf.StsDurationSeconds,
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (r *s3Repo) ListBuckets(ctx context.Context) ([]string, error) {
	output, err := r.data.s3.Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		r.logger.Error("list buckets failed", "error", err)
		return nil, err
	}
	buckets := make([]string, len(output.Buckets))
	for i, bucket := range output.Buckets {
		buckets[i] = *bucket.Name
	}
	return buckets, nil
}
