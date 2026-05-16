package service

import (
	"context"

	pb "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UploadService struct {
	pb.UnimplementedUploadServiceServer

	logger *log.Helper
	uc     *biz.S3Usecase
}

func NewUploadService(logger log.Logger, uc *biz.S3Usecase) *UploadService {
	return &UploadService{
		logger: log.NewHelper(logger),
		uc:     uc,
	}
}

// GetUploadSession 获取上传会话
func (s *UploadService) GetUploadPresignedURL(ctx context.Context, req *emptypb.Empty) (*pb.GetUploadPresignedURLReply, error) {
	url, err := s.uc.GetUploadPresignedURL(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetUploadPresignedURLReply{
		Url: url,
	}, nil
}

func (s *UploadService) GetUploadSession(ctx context.Context, req *emptypb.Empty) (*pb.GetUploadSessionReply, error) {
	session, err := s.uc.GetUploadSession(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetUploadSessionReply{
		AccessKeyId:     session.AccessKeyID,
		SecretAccessKey: session.SecretAccessKey,
		SessionToken:    session.SessionToken,
	}, nil
}

func (s *UploadService) ListBuckets(ctx context.Context, req *emptypb.Empty) (*pb.ListBucketsReply, error) {
	buckets, err := s.uc.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.ListBucketsReply{
		Buckets: buckets,
	}, nil
}
