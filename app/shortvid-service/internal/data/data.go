package data

import (
	"shortvid-backend/app/shortvid-service/internal/data/infra/cache"
	"shortvid-backend/app/shortvid-service/internal/data/infra/db"
	"shortvid-backend/app/shortvid-service/internal/data/infra/storage"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(db.NewDB, cache.NewRedis, storage.NewMinioClient, NewData, NewTxRepo, NewAccountRepo, NewUserRepo, NewUserSessionRepo)

// 基础设施的数据模型
type Data struct {
	db     *gorm.DB
	redis  *redis.Client
	minio  *minio.Client
	logger log.Logger
}

// NewData .
func NewData(db *gorm.DB, redis *redis.Client, minio *minio.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Log(log.LevelInfo, "msg", "closing the data resources")
	}
	return &Data{
		db:     db,
		redis:  redis,
		minio:  minio,
		logger: logger,
	}, cleanup, nil
}
