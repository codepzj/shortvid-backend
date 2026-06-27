package data

import (
	"context"
	"errors"
	"shortvid-backend/app/user-service/internal/biz"
	"shortvid-backend/app/user-service/internal/data/model"

	"gorm.io/gorm"
)

type accountRepo struct {
	data *Data
}

func NewAccountRepo(data *Data) biz.AccountRepo {
	return &accountRepo{data: data}
}

func (r *accountRepo) CreateAccountWithTx(ctx context.Context, tx *gorm.DB, account *model.Account) error {
	return tx.WithContext(ctx).Create(account).Error
}

func (r *accountRepo) GetByEmailAndProvider(ctx context.Context, email string, provider string) (*model.Account, error) {
	var account model.Account
	err := r.data.db.WithContext(ctx).Where("email = ? AND provider = ?", email, provider).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}
