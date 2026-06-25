package biz

import (
	"context"
	"shortvid-backend/app/user-service/internal/data/model"

	"gorm.io/gorm"
)

type AccountRepo interface {
	CreateAccountWithTx(ctx context.Context, tx *gorm.DB, account *model.Account) error
	GetByEmailAndProvider(ctx context.Context, email string, provider string) (*model.Account, error)
}

type AccountUsecase struct {
	repo AccountRepo
}

func NewAccountUsecase(repo AccountRepo) *AccountUsecase {
	return &AccountUsecase{repo: repo}
}
