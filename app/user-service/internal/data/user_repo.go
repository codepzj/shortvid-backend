package data

import (
	"context"
	"errors"
	"shortvid-backend/app/user-service/internal/biz"
	"shortvid-backend/app/user-service/internal/data/model"
	"time"

	"gorm.io/gorm"
)

type userRepo struct {
	data *Data
}

func NewUserRepo(data *Data) biz.UsersRepo {
	return &userRepo{data: data}
}

func (r *userRepo) CreateUser(ctx context.Context, user *model.User) error {
	return r.data.db.WithContext(ctx).Create(user).Error
}

func (r *userRepo) CreateUserWithTx(ctx context.Context, tx *gorm.DB, user *model.User) error {
	return tx.WithContext(ctx).Create(user).Error
}

func (r *userRepo) GetUserByEmailAndProvider(ctx context.Context, email string, provider string) (*model.User, error) {
	var user model.User
	err := r.data.db.WithContext(ctx).Where("email = ? AND provider = ?", email, provider).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetUserByUID(ctx context.Context, UID int) (*model.User, error) {
	var user model.User
	err := r.data.db.WithContext(ctx).Where("uid = ?", UID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) UpdateLoginInfo(ctx context.Context, UID int) error {
	now := time.Now()
	return r.data.db.WithContext(ctx).Model(&model.User{}).Where("uid = ?", UID).Updates(map[string]any{
		"last_login_at": now,
		"login_count":   gorm.Expr("login_count + 1"),
		"updated_at":    now,
	}).Error
}
