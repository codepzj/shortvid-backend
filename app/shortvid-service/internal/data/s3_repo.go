package data

import (
	"context"
	"fmt"
	"shortvid-backend/app/shortvid-service/internal/biz"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/go-kratos/kratos/v2/log"
)

type s3Repo struct {
	data   *Data
	logger *log.Helper
}

func NewS3Repo(data *Data, logger log.Logger) biz.S3Repo {
	return &s3Repo{data: data, logger: log.NewHelper(logger)}
}

func (r *s3Repo) GetUploadSession(ctx context.Context) (*sts.GetSessionTokenOutput, error) {
	output, err := r.data.s3.StsClient.GetSessionToken(ctx, &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int32(7200),
	})
	if err != nil {
		r.logger.Error("get upload session failed")
		return nil, err
	}

	fmt.Println(output.Credentials.AccessKeyId)
	fmt.Println(output.Credentials.SecretAccessKey)
	fmt.Println(output.Credentials.SessionToken)
	fmt.Println(output.Credentials.Expiration)
	return output, nil
}
