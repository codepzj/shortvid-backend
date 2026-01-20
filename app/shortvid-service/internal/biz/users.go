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
	logger *log.Helper
	repo   UsersRepo
}

func NewUsersUsecase(logger log.Logger, repo UsersRepo) *UsersUsecase {
	return &UsersUsecase{logger: log.NewHelper(logger), repo: repo}
}

func (uc *UsersUsecase) FindOrCreateUser(ctx context.Context, user *User) (*User, bool, error) {
	existingUser, err := uc.repo.GetUserByEmailAndProvider(ctx, user.Email, user.Provider)
	if err != nil {
		return nil, false, err
	}
	if existingUser != nil {
		return existingUser, false, nil
	}
	err = uc.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, false, err
	}
	return user, true, nil
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
