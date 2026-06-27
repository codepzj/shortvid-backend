package data

import (
	"shortvid-backend/app/user-service/internal/data/infra/cache"
	"shortvid-backend/app/user-service/internal/data/infra/db"

	"github.com/go-kratos/kratos/v3/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(db.NewDB, cache.NewRedis, NewData, NewTxRepo, NewAccountRepo, NewUserRepo, NewUserSessionRepo)

// 基础设施的数据模型
type Data struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewData 初始化基础设施
func NewData(db *gorm.DB, redis *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the infra data resources")
	}
	return &Data{
		db:    db,
		redis: redis,
	}, cleanup, nil
}
