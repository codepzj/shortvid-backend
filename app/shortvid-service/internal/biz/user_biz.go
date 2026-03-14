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
	"gorm.io/gorm"
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
	UID         int
	Ctime       int64  // 创建时间
	Nickname    string // 昵称
	Avatar      string // 头像
	Email       string // 邮箱
	Provider    string // 提供商
	ProviderUID string // 提供商UID
}

type UsersRepo interface {
	CreateUser(ctx context.Context, user *model.User) error
	CreateUserWithTx(ctx context.Context, tx *gorm.DB, user *model.User) error
	GetUserByEmailAndProvider(ctx context.Context, email string, provider string) (*model.User, error)
	GetUserByUID(ctx context.Context, UID int) (*model.User, error)
	UpdateLoginInfo(ctx context.Context, userUID int) error
}

type UsersUsecase struct {
	logger      log.Logger
	repo        UsersRepo
	accountRepo AccountRepo
}

func NewUsersUsecase(logger log.Logger, repo UsersRepo, accountRepo AccountRepo) *UsersUsecase {
	return &UsersUsecase{logger: logger, repo: repo, accountRepo: accountRepo}
}

// FirebaseFindOrCreateUser 查询或创建用户[Firebase]
func (uc *UsersUsecase) FirebaseFindOrCreateUser(ctx context.Context, user *User) (*UserProfile, bool, error) {
	// 判断账户是否存在
	existAcc, err := uc.accountRepo.GetByEmailAndProvider(ctx, user.Email, user.Provider)
	if err != nil {
		uc.logger.Log(log.LevelError, "msg", "Get account by email and provider failed", "error", err)
		return nil, false, err
	}
	if existAcc != nil {
		// 账户存在, 继续查找用户
		existUser, err := uc.repo.GetUserByUID(ctx, existAcc.UID)
		if err != nil {
			uc.logger.Log(log.LevelError, "msg", "Get user by uid failed", "error", err)
			return nil, false, err
		}
		if existUser == nil {
			return nil, false, errors.New("account exists, but user not found")
		}

		return &UserProfile{
			ID:          existUser.ID,
			UID:         existUser.UID,
			Ctime:       existUser.CreatedAt.Unix(),
			Nickname:    existUser.Nickname,
			Avatar:      existUser.Avatar,
			Email:       existAcc.Email,
			Provider:    existAcc.Provider,
			ProviderUID: existAcc.ProviderUID,
		}, false, nil
	}

	// 账户不存在, 创建账户和用户
	maxCount := 10
	userModel := &model.User{
		UID:         0,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		LastLoginAt: time.Now(),
	}
	for range maxCount {
		// 生成唯一的UserUID
		userModel.UID = uc.generateUniqueUserUID()
		err := uc.repo.CreateUser(ctx, userModel)

		if err != nil {
			var mysqlErr *mysqlDriver.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				uc.logger.Log(log.LevelInfo, "msg", "UserUID already exists, retrying", "userUID", userModel.UID)
				continue
			}
			uc.logger.Log(log.LevelWarn, "msg", "Create user failed", "error", err)
			return nil, false, err
		}
		return &UserProfile{
			ID:       userModel.ID,
			UID:      userModel.UID,
			Nickname: userModel.Nickname,
			Avatar:   userModel.Avatar,
		}, true, nil
	}
	return nil, false, errors.New("create user failed after max retries")
}

// GetUserByUserUID 根据UserUID查询用户
func (uc *UsersUsecase) GetUserByUID(ctx context.Context, UID int) (*UserProfile, error) {
	userModel, err := uc.repo.GetUserByUID(ctx, UID)
	if err != nil {
		uc.logger.Log(log.LevelError, "msg", "Get user by userUid failed", "error", err)
		return nil, err
	}
	if userModel == nil {
		return nil, v1.ErrorUserNotFound("user not found")
	}
	return &UserProfile{
		ID:       userModel.ID,
		UID:      userModel.UID,
		Nickname: userModel.Nickname,
		Avatar:   userModel.Avatar,
	}, nil
}

// UpdateLoginInfo 更新登录信息
func (uc *UsersUsecase) UpdateLoginInfo(ctx context.Context, userUID int) error {
	return uc.repo.UpdateLoginInfo(ctx, userUID)
}

// generateUniqueUserUID 生成唯一的UserUID (10000-999999999范围)
func (uc *UsersUsecase) generateUniqueUserUID() int {
	const minUID = 10000
	const maxUID = 999999999
	return rand.Intn(maxUID-minUID+1) + minUID
}
