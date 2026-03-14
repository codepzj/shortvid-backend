package biz

import (
	"context"
	"shortvid-backend/app/shortvid-service/internal/data/model"
)

type AccountRepo interface {
	CreateAccount(ctx context.Context, account *model.Account) error
	GetByEmailAndProvider(ctx context.Context, email string, provider string) (*model.Account, error)
}

type AccountUsecase struct {
	repo AccountRepo
}

func NewAccountUsecase(repo AccountRepo) *AccountUsecase {
	return &AccountUsecase{repo: repo}
}
