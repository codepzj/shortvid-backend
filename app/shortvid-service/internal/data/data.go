package data

import (
	"shortvid-backend/app/shortvid-service/internal/data/infra/cache"
	"shortvid-backend/app/shortvid-service/internal/data/infra/db"
	"shortvid-backend/app/shortvid-service/internal/data/infra/storage"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(db.NewDB, cache.NewRedis, storage.NewS3, NewData, NewTxRepo, NewAccountRepo, NewUserRepo, NewUserSessionRepo)

// 基础设施的数据模型
type Data struct {
	db     *gorm.DB
	redis  *redis.Client
	s3     *storage.S3Data
	logger *log.Helper
}

// NewData 初始化基础设施
func NewData(db *gorm.DB, redis *redis.Client, s3 *storage.S3Data, logger log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(logger)
	cleanup := func() {
		helper.Infow("msg", "closing the infra data resources")
	}
	return &Data{
		db:     db,
		redis:  redis,
		s3:     s3,
		logger: helper,
	}, cleanup, nil
}
