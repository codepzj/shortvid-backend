package data

import (
	"context"
	"shortvid-backend/app/shortvid-service/internal/biz"
	"shortvid-backend/app/shortvid-service/internal/conf"
	"shortvid-backend/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/go-kratos/kratos/v3/log"
)

type s3Repo struct {
	conf *conf.S3
	data *Data
}

func NewS3Repo(conf *conf.S3, data *Data) biz.S3Repo {
	return &s3Repo{
		conf: conf,
		data: data,
	}
}

func (r *s3Repo) GetUploadSession(ctx context.Context, vgroup string) (*sts.AssumeRoleOutput, string, string, error) {
	// resource: arn:aws:s3:::bucket_name/prefix/*
	path := utils.GetVidPoolPathFromVgroup(vgroup)
	resource := "arn:aws:s3:::" + r.conf.Bucket + "/" + path
	stsPolicy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"s3:PutObject"
				],
				"Resource": "` + resource + `"
			}
		]
	}`
	output, err := r.data.s3.StsClient.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         &r.conf.StsRoleArn,
		RoleSessionName: &r.conf.StsSessionName,
		Policy:          &stsPolicy,
		DurationSeconds: &r.conf.StsDurationSeconds,
	})
	if err != nil {
		return nil, "", "", err
	}
	return output, r.conf.Bucket, path, nil
}

func (r *s3Repo) ListBuckets(ctx context.Context) ([]string, error) {
	output, err := r.data.s3.Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		log.ErrorContext(ctx, "list buckets failed", "error", err)
		return nil, err
	}
	buckets := make([]string, len(output.Buckets))
	for i, bucket := range output.Buckets {
		buckets[i] = *bucket.Name
	}
	return buckets, nil
}
