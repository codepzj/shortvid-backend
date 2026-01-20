package biz

import (
	"context"
	v1 "shortvid-backend/api/shortvid-service/v1"

	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	UserUID     string
	Nickname    string
	Avatar      string
	Email       string
	ProviderUID string
	Provider    string
}

type UsersRepo interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id int32) (*User, error)
	GetUserByEmailAndProvider(ctx context.Context, email string, provider string) (*User, error)
}

type UsersUsecase struct {
	logger log.Logger
	repo   UsersRepo
}

func NewUsersUsecase(logger log.Logger, repo UsersRepo) *UsersUsecase {
	return &UsersUsecase{logger: logger, repo: repo}
}

func (uc *UsersUsecase) FindOrCreateUser(ctx context.Context, user *User) (*User, bool, error) {
	existingUser, err := uc.repo.GetUserByEmailAndProvider(ctx, user.Email, user.Provider)
	if err != nil {
		uc.logger.Log(log.LevelError, "msg", "Get user by email and provider failed", "error", err)
		return nil, false, err
	}
	if existingUser != nil {
		return existingUser, false, nil
	}
	err = uc.repo.CreateUser(ctx, user)
	if err != nil {
		uc.logger.Log(log.LevelError, "msg", "Create user failed", "error", err)
		return nil, false, err
	}
	return user, true, nil
}

func (uc *UsersUsecase) GetUserByID(ctx context.Context, id int32) (*User, error) {
	user, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		uc.logger.Log(log.LevelError, "msg", "Get user by id failed", "error", err)
		return nil, err
	}
	if user == nil {
		return nil, v1.ErrorUserNotFound("user not found")
	}
	return user, nil
}
