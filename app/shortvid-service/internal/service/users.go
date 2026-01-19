package service

import (
	"context"
	v1 "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/biz"
)

type UsersService struct {
	v1.UnimplementedUsersServiceServer

	uc *biz.UsersUsecase
}

func NewUsersService(uc *biz.UsersUsecase) *UsersService {
	return &UsersService{uc: uc}
}

func (s *UsersService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	err := s.uc.CreateUser(ctx, &biz.User{
		Nickname:    req.Nickname,
		Avatar:      req.Avatar,
		Email:       req.Email,
		Provider:    req.Provider,
		ProviderUID: req.ProviderUid,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateUserResponse{}, nil
}

func (s *UsersService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	user, err := s.uc.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserResponse{
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Email:       user.Email,
		Provider:    user.Provider,
		ProviderUid: user.ProviderUID,
	}, nil
}
