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

func (r *usersRepo) CreateUser(ctx context.Context, user *model.User) error {
	return r.data.query.User.WithContext(ctx).Create(user)
}

func (r *usersRepo) GetUserByID(ctx context.Context, id int32) (*model.User, error) {
	user, err := r.data.query.User.WithContext(ctx).Where(r.data.query.User.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *usersRepo) GetUserByEmailAndProvider(ctx context.Context, email string, provider string) (*model.User, error) {
	user, err := r.data.query.User.WithContext(ctx).Where(r.data.query.User.Email.Eq(email), r.data.query.User.Provider.Eq(provider)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
