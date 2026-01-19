package data

import (
	"context"
	"errors"
	"shortvid-backend/app/shortvid-service/internal/biz"
	"shortvid-backend/app/shortvid-service/internal/data/model"

	"gorm.io/gorm"
)

type usersRepo struct {
	data *Data
}

func NewUsersRepo(data *Data) biz.UsersRepo {
	return &usersRepo{data: data}
}

func (r *usersRepo) CreateUser(ctx context.Context, user *biz.User) error {
	return r.data.query.User.WithContext(ctx).Create(&model.User{
		Nickname:    user.Nickname,
		Avatar:      &user.Avatar,
		Email:       &user.Email,
		ProviderUID: &user.ProviderUID,
		Provider:    &user.Provider,
	})
}

func (r *usersRepo) GetUserByID(ctx context.Context, id int32) (*biz.User, error) {
	user, err := r.data.query.User.WithContext(ctx).Where(r.data.query.User.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &biz.User{
		Nickname:    user.Nickname,
		Avatar:      *user.Avatar,
		Email:       *user.Email,
		ProviderUID: *user.ProviderUID,
		Provider:    *user.Provider,
	}, nil
}
