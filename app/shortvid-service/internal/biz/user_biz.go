package biz

import (
	"context"
	"errors"
	"math/rand"
	v1 "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/data/model"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserDTO struct {
	Nickname    string
	Avatar      string
	Email       string
	ProviderUID string
	Provider    string
}

type UserProfileVO struct {
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
	txRepo      TxRepo
	repo        UsersRepo
	accountRepo AccountRepo
}

func NewUsersUsecase(logger log.Logger, txRepo TxRepo, repo UsersRepo, accountRepo AccountRepo) *UsersUsecase {
	return &UsersUsecase{logger: logger, txRepo: txRepo, repo: repo, accountRepo: accountRepo}
}

// FirebaseFindOrCreateUser 查询或创建用户[Firebase]
func (uc *UsersUsecase) FirebaseFindOrCreateUser(ctx context.Context, dto *UserDTO) (*UserProfileVO, bool, error) {
	// 判断账户是否存在
	existAcc, err := uc.accountRepo.GetByEmailAndProvider(ctx, dto.Email, dto.Provider)
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

		return &UserProfileVO{
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

	// 创建账户和用户(强关联操作)
	user := &model.User{
		Nickname:    dto.Nickname,
		Avatar:      dto.Avatar,
		LastLoginAt: time.Now(),
	}
	account := &model.Account{
		Email:       dto.Email,
		Provider:    dto.Provider,
		ProviderUID: dto.ProviderUID,
	}
	err = uc.txRepo.ExecFunc(func(tx *gorm.DB) error {

		maxCount := 10
		for range maxCount {
			// 生成唯一的UID
			user.UID = uc.generateUniqueUserUID()
			err := uc.repo.CreateUserWithTx(ctx, tx, user)

			if err != nil {
				var mysqlErr *mysql.MySQLError
				if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
					uc.logger.Log(log.LevelInfo, "msg", "UserUID already exists, retrying", "userUID", user.UID)
					continue
				}
				uc.logger.Log(log.LevelError, "msg", "Create user failed", "error", err)
				return err
			}
			break
		}

		account.UID = user.UID
		err := uc.accountRepo.CreateAccountWithTx(ctx, tx, account)
		if err != nil {
			uc.logger.Log(log.LevelError, "msg", "Create account failed", "error", err)
			return err
		}

		return nil
	})

	if err != nil {
		uc.logger.Log(log.LevelError, "msg", "Create user and account failed", "error", err)
		return nil, false, err
	}

	return &UserProfileVO{
		ID:          user.ID,
		UID:         user.UID,
		Ctime:       user.CreatedAt.Unix(), // 创建时间
		Nickname:    user.Nickname,         // 昵称
		Avatar:      user.Avatar,           // 头像
		Email:       account.Email,         // 邮箱
		Provider:    account.Provider,      // 提供商
		ProviderUID: account.ProviderUID,   // 提供商UID
	}, true, nil
}

// GetUserByUserUID 根据UserUID查询用户
func (uc *UsersUsecase) GetUserByUID(ctx context.Context, UID int) (*UserProfileVO, error) {
	user, err := uc.repo.GetUserByUID(ctx, UID)
	if err != nil {
		uc.logger.Log(log.LevelError, "msg", "Get user by userUid failed", "error", err)
		return nil, err
	}
	if user == nil {
		return nil, v1.ErrorUserNotFound("user not found")
	}
	return &UserProfileVO{
		ID:       user.ID,
		UID:      user.UID,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
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
