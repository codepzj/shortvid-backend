package service

import (
	"context"

	uploadV1 "shortvid-backend/api/v1/upload"
	"shortvid-backend/app/shortvid-service/internal/biz"

	"google.golang.org/protobuf/types/known/emptypb"
)

type UploadService struct {
	uploadV1.UnimplementedUploadServiceServer

	uc *biz.S3Usecase
}

func NewUploadService(uc *biz.S3Usecase) *UploadService {
	return &UploadService{
		uc: uc,
	}
}

func (s *UploadService) GetUploadSession(ctx context.Context, req *emptypb.Empty) (*uploadV1.GetUploadSessionReply, error) {
	session, err := s.uc.GetUploadSession(ctx)
	if err != nil {
		return nil, err
	}
	return &uploadV1.GetUploadSessionReply{
		AccessKeyId:     session.AccessKeyID,
		SecretAccessKey: session.SecretAccessKey,
		SessionToken:    session.SessionToken,
	}, nil
}

func (s *UploadService) ListBuckets(ctx context.Context, req *emptypb.Empty) (*uploadV1.ListBucketsReply, error) {
	buckets, err := s.uc.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}
	return &uploadV1.ListBucketsReply{
		Buckets: buckets,
	}, nil
}
