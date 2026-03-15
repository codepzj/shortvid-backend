package data

import (
	"shortvid-backend/app/shortvid-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type TxRepo struct {
	data *Data
}

func NewTxRepo(data *Data) biz.TxRepo {
	return &TxRepo{
		data: data,
	}
}

func (r *TxRepo) ExecFunc(fn func(tx *gorm.DB) error) error {
	err := r.data.db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
	if err != nil {
		r.data.logger.Log(log.LevelError, "error", "transaction failed")
		return err
	}
	return nil
}
