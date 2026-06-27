package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	healthV1 "shortvid-backend/api/v1/health"
)

type HealthService struct {
	healthV1.UnimplementedHealthServiceServer
}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) Health(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
