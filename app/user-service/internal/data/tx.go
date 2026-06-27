package data

import (
	"log/slog"
	"shortvid-backend/app/user-service/internal/biz"

	"github.com/go-kratos/kratos/v3/log"
	"gorm.io/gorm"
)

type txRepo struct {
	data *Data
}

func NewTxRepo(data *Data) biz.TxRepo {
	return &txRepo{
		data: data,
	}
}

func (r *txRepo) ExecFunc(fn func(tx *gorm.DB) error) error {
	err := r.data.db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
	if err != nil {
		log.Error("transaction failed", slog.Any("error", err))
		return err
	}
	return nil
}
