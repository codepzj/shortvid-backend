package data

import (
	"shortvid-backend/app/shortvid-service/internal/data/infra/cache"
	"shortvid-backend/app/shortvid-service/internal/data/infra/db"
	"shortvid-backend/app/shortvid-service/internal/data/query"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(db.NewDB, cache.NewRedis, NewData, NewUsersRepo, NewUserSessionRepo)

// Data .
type Data struct {
	query  *query.Query
	redis  *redis.Client
	logger log.Logger
}

// NewData .
func NewData(db *gorm.DB, redis *redis.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Log(log.LevelInfo, "closing the data resources")
	}
	return &Data{
		query:  query.Use(db),
		redis:  redis,
		logger: logger,
	}, cleanup, nil
}
