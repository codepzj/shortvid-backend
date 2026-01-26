package biz

import (
	"context"
	"errors"
	"math/rand"
	v1 "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/data/model"
	"time"

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
	ID          int
	UserUID     int
	Nickname    string
	Avatar      string
	Email       string
	Provider    string
	ProviderUID string
}

type UsersRepo interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByEmailAndProvider(ctx context.Context, email string, provider string) (*model.User, error)
	GetUserByUserUID(ctx context.Context, userUID int) (*model.User, error)
	UpdateLoginInfo(ctx context.Context, userID int) error
}

type UsersUsecase struct {
	logger log.Logger
	repo   UsersRepo
}

func NewUsersUsecase(logger log.Logger, repo UsersRepo) *UsersUsecase {
	return &UsersUsecase{logger: logger, repo: repo}
}

// FindOrCreateUser 查询或创建用户
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
			Avatar:      existingUserModel.Avatar,
			Email:       existingUserModel.Email,
			ProviderUID: existingUserModel.ProviderUID,
			Provider:    existingUserModel.Provider,
		}, false, nil
	}

	// 生成唯一的UserUID
	userUID, err := uc.generateUniqueUserUID(ctx)
	if err != nil {
		uc.logger.Log(log.LevelError, "msg", "Generate unique UserUID failed", "error", err)
		return nil, false, err
	}

	userModel := &model.User{
		UserUID:     userUID,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Email:       user.Email,
		ProviderUID: user.ProviderUID,
		Provider:    user.Provider,
		LastLoginAt: time.Now(),
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
		Avatar:      userModel.Avatar,
		Email:       userModel.Email,
		ProviderUID: userModel.ProviderUID,
		Provider:    userModel.Provider,
	}, true, nil
}

// GetUserByID 根据ID查询用户
func (uc *UsersUsecase) GetUserByID(ctx context.Context, id int) (*UserProfile, error) {
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
		Avatar:      userModel.Avatar,
		Email:       userModel.Email,
		Provider:    userModel.Provider,
		ProviderUID: userModel.ProviderUID,
	}, nil
}

// GetUserByUID 根据UID查询用户
func (uc *UsersUsecase) GetUserByUID(ctx context.Context, uid int) (*UserProfile, error) {
	userModel, err := uc.repo.GetUserByUserUID(ctx, uid)
	if err != nil {
		uc.logger.Log(log.LevelError, "msg", "Get user by uid failed", "error", err)
		return nil, err
	}
	if userModel == nil {
		return nil, v1.ErrorUserNotFound("user not found")
	}
	return &UserProfile{
		ID:          userModel.ID,
		UserUID:     userModel.UserUID,
		Nickname:    userModel.Nickname,
		Avatar:      userModel.Avatar,
		Email:       userModel.Email,
		Provider:    userModel.Provider,
		ProviderUID: userModel.ProviderUID,
	}, nil
}

// UpdateLoginInfo 更新登录信息
func (uc *UsersUsecase) UpdateLoginInfo(ctx context.Context, userID int) error {
	return uc.repo.UpdateLoginInfo(ctx, userID)
}

// generateUniqueUserUID 生成唯一的UserUID (10000-999999999范围)
func (uc *UsersUsecase) generateUniqueUserUID(ctx context.Context) (int, error) {
	const maxRetries = 10
	const minUID = 10000
	const maxUID = 999999999

	for i := range maxRetries {
		// 生成随机UserUID
		userUID := rand.Intn(maxUID-minUID+1) + minUID

		// 检查是否已存在
		existingUser, err := uc.repo.GetUserByUserUID(ctx, userUID)
		if err != nil {
			return 0, err
		}

		// 如果不存在，返回这个UserUID
		if existingUser == nil {
			return userUID, nil
		}

		uc.logger.Log(log.LevelWarn, "msg", "Generated UserUID already exists, retrying", "userUID", userUID, "attempt", i+1)
	}

	return 0, errors.New("failed to generate unique UserUID after max retries")
}
