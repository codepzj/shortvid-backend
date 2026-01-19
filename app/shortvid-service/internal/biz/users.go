package biz

import (
	"context"
	v1 "shortvid-backend/api/shortvid-service/v1"

	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	Nickname    string
	Avatar      string
	Email       string
	ProviderUID string
	Provider    string
}

type UsersRepo interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id int32) (*User, error)
}

type UsersUsecase struct {
	logger *log.Helper
	repo   UsersRepo
}

func NewUsersUsecase(logger log.Logger, repo UsersRepo) *UsersUsecase {
	return &UsersUsecase{logger: log.NewHelper(logger), repo: repo}
}

func (uc *UsersUsecase) CreateUser(ctx context.Context, user *User) error {
	return uc.repo.CreateUser(ctx, user)
}

func (uc *UsersUsecase) GetUserByID(ctx context.Context, id int32) (*User, error) {
	user, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		uc.logger.Errorf("get user by id failed: %v", err)
		return nil, err
	}
	if user == nil {
		return nil, v1.ErrorUserNotFound("user not found")
	}
	return user, nil
}
