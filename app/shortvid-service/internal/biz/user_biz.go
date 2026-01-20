package biz

import (
	"context"
	v1 "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	Nickname    string
	Avatar      string
	Email       string
	ProviderUID string
	Provider    string
}

type UserProfile struct {
	ID          int32
	UserUID     int32
	Nickname    string
	Avatar      string
	Email       string
	Provider    string
	ProviderUID string
}

type UsersRepo interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int32) (*model.User, error)
	GetUserByEmailAndProvider(ctx context.Context, email string, provider string) (*model.User, error)
}

type UsersUsecase struct {
	logger log.Logger
	repo   UsersRepo
}

func NewUsersUsecase(logger log.Logger, repo UsersRepo) *UsersUsecase {
	return &UsersUsecase{logger: logger, repo: repo}
}

func (uc *UsersUsecase) FindOrCreateUser(ctx context.Context, user *User) (*UserProfile, bool, error) {
	existingUserModel, err := uc.repo.GetUserByEmailAndProvider(ctx, user.Email, user.Provider)
	if err != nil {
		uc.logger.Log(log.LevelError, "msg", "Get user by email and provider failed", "error", err)
		return nil, false, err
	}
	if existingUserModel != nil {
		return &UserProfile{
			ID:          existingUserModel.ID,
			UserUID:     existingUserModel.UserUID,
			Nickname:    existingUserModel.Nickname,
			Avatar:      *existingUserModel.Avatar,
			Email:       *existingUserModel.Email,
			ProviderUID: *existingUserModel.ProviderUID,
			Provider:    *existingUserModel.Provider,
		}, false, nil
	}
	userModel := &model.User{
		Nickname:    user.Nickname,
		Avatar:      &user.Avatar,
		Email:       &user.Email,
		ProviderUID: &user.ProviderUID,
		Provider:    &user.Provider,
	}
	err = uc.repo.CreateUser(ctx, userModel)
	if err != nil {
		uc.logger.Log(log.LevelError, "msg", "Create user failed", "error", err)
		return nil, false, err
	}
	return &UserProfile{
		ID:          userModel.ID,
		UserUID:     userModel.UserUID,
		Nickname:    userModel.Nickname,
		Avatar:      *userModel.Avatar,
		Email:       *userModel.Email,
		ProviderUID: *userModel.ProviderUID,
		Provider:    *userModel.Provider,
	}, true, nil
}

func (uc *UsersUsecase) GetUserByID(ctx context.Context, id int32) (*UserProfile, error) {
	userModel, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		uc.logger.Log(log.LevelError, "msg", "Get user by id failed", "error", err)
		return nil, err
	}
	if userModel == nil {
		return nil, v1.ErrorUserNotFound("user not found")
	}
	return &UserProfile{
		ID:          userModel.ID,
		UserUID:     userModel.UserUID,
		Nickname:    userModel.Nickname,
		Avatar:      *userModel.Avatar,
		Email:       *userModel.Email,
		Provider:    *userModel.Provider,
		ProviderUID: *userModel.ProviderUID,
	}, nil
}
