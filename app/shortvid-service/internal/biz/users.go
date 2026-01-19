package biz

import (
	"context"
)

type User struct {
	Nickname   string
	Avatar     string
	Email      string
	ProviderUID string
	Provider   string
}

type UsersRepo interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id int32) (*User, error)
}

type UsersUsecase struct {
	repo UsersRepo
}

func NewUsersUsecase(repo UsersRepo) *UsersUsecase {
	return &UsersUsecase{repo: repo}
}

func (uc *UsersUsecase) CreateUser(ctx context.Context, user *User) error {
	return uc.repo.CreateUser(ctx, user)
}

func (uc *UsersUsecase) GetUserByID(ctx context.Context, id int32) (*User, error) {
	return uc.repo.GetUserByID(ctx, id)
}