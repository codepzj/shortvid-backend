package service

import (
	"context"

	pb "shortvid-backend/api/shortvid-service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type HealthService struct {
	pb.UnimplementedHealthServiceServer
}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) Health(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
    return &emptypb.Empty{}, nil
}
