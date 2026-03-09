package data

import (
	"context"
	"errors"
	"shortvid-backend/app/shortvid-service/internal/biz"
	"shortvid-backend/app/shortvid-service/internal/data/model"
	"time"

	"gorm.io/gorm"
)

type usersRepo struct {
	data *Data
}

func NewUsersRepo(data *Data) biz.UsersRepo {
	return &usersRepo{data: data}
}

func (r *usersRepo) CreateUser(ctx context.Context, user *model.User) error {
	return r.data.db.WithContext(ctx).Create(user).Error
}

func (r *usersRepo) CreateUserWithTx(ctx context.Context, tx *gorm.DB, user *model.User) error {
	return tx.WithContext(ctx).Create(user).Error
}

func (r *usersRepo) GetUserByEmailAndProvider(ctx context.Context, email string, provider string) (*model.User, error) {
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

func (r *usersRepo) GetUserByUserUID(ctx context.Context, userUID int) (*model.User, error) {
	var user model.User
	err := r.data.db.WithContext(ctx).Where("user_uid = ?", userUID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *usersRepo) UpdateLoginInfo(ctx context.Context, userUID int) error {
	now := time.Now()
	return r.data.db.WithContext(ctx).Model(&model.User{}).Where("user_uid = ?", userUID).Updates(map[string]any{
		"last_login_at": now,
		"login_count":   gorm.Expr("login_count + 1"),
		"updated_at":    now,
	}).Error
}
