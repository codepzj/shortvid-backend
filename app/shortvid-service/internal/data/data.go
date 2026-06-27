package data

import (
	"shortvid-backend/app/shortvid-service/internal/data/infra/db"
	"shortvid-backend/app/shortvid-service/internal/data/infra/storage"

	"github.com/go-kratos/kratos/v3/log"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(db.NewDB, storage.NewS3, NewData, NewS3Repo)

// 基础设施的数据模型
type Data struct {
	db *gorm.DB
	s3 *storage.S3Data
}

// NewData 初始化基础设施
func NewData(db *gorm.DB, s3 *storage.S3Data) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the infra data resources")
	}
	return &Data{
		db: db,
		s3: s3,
	}, cleanup, nil
}
