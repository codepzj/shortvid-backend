package data

import (
	"context"
	"shortvid-backend/app/shortvid-service/internal/biz"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/go-kratos/kratos/v2/log"
)

type s3Repo struct {
	data   *Data
	logger *log.Helper
}

func NewS3Repo(data *Data, logger log.Logger) biz.S3Repo {
	return &s3Repo{
		data:   data,
		logger: log.NewHelper(logger),
	}
}

func (r *s3Repo) GetUploadPresignedURL(ctx context.Context) (string, error) {
	output, err := r.data.s3.PresignedClient.PresignGetObject(ctx,
		&s3.GetObjectInput{
			Bucket: aws.String("test"),
			Key:    aws.String("test.txt"),
		},
		s3.WithPresignExpires(7200*time.Second),
	)
	if err != nil {
		return "", err
	}
	return output.URL, nil
}

func (r *s3Repo) GetUploadSession(ctx context.Context) (*sts.AssumeRoleOutput, error) {
	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"s3:GetObject",
					"s3:PutObject",
					"s3:GetBucketLocation",
					"s3:DeleteObject"
				],
				"Resource": [
					"arn:aws:s3:::*"
				]
			}
		]
	}`
	output, err := r.data.s3.StsClient.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         aws.String("arn:aws:iam::codepzj:role/readwrite"),
		RoleSessionName: aws.String(""),
		Policy:          aws.String(policy),
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
