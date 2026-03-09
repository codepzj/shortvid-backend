package biz

import (
	"context"
	"errors"
	"math/rand"
	v1 "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/data/model"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	mysqlDriver "github.com/go-sql-driver/mysql"
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
		}, false, nil
	}

	maxCount := 10
	userModel := &model.User{
		UserUID:     0,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		LastLoginAt: time.Now(),
	}
	for range maxCount {
		// 生成唯一的UserUID
		userModel.UserUID = uc.generateUniqueUserUID()
		err := uc.repo.CreateUser(ctx, userModel)

		if err != nil {
			var mysqlErr *mysqlDriver.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				uc.logger.Log(log.LevelInfo, "msg", "UserUID already exists, retrying", "userUID", userModel.UserUID)
				continue
			}
			uc.logger.Log(log.LevelError, "msg", "Create user failed", "error", err)
			return nil, false, err
		}
		return &UserProfile{
			ID:          userModel.ID,
			UserUID:     userModel.UserUID,
			Nickname:    userModel.Nickname,
			Avatar:      userModel.Avatar,
		}, true, nil
	}
	return nil, false, errors.New("create user failed after max retries")
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
	}, nil
}

// GetUserByUserUID 根据UserUID查询用户
func (uc *UsersUsecase) GetUserByUserUID(ctx context.Context, userUID int) (*UserProfile, error) {
	userModel, err := uc.repo.GetUserByUserUID(ctx, userUID)
	if err != nil {
		uc.logger.Log(log.LevelError, "msg", "Get user by userUid failed", "error", err)
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
	}, nil
}

// UpdateLoginInfo 更新登录信息
func (uc *UsersUsecase) UpdateLoginInfo(ctx context.Context, userID int) error {
	return uc.repo.UpdateLoginInfo(ctx, userID)
}

// generateUniqueUserUID 生成唯一的UserUID (10000-999999999范围)
func (uc *UsersUsecase) generateUniqueUserUID() int {
	const minUID = 10000
	const maxUID = 999999999
	return rand.Intn(maxUID-minUID+1) + minUID
}
