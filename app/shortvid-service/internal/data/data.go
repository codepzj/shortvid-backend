package data

import (
	"shortvid-backend/app/shortvid-service/internal/data/infra"
	"shortvid-backend/app/shortvid-service/internal/data/infra/cache"
	"shortvid-backend/app/shortvid-service/internal/data/infra/db"
	"shortvid-backend/app/shortvid-service/internal/data/query"

	"firebase.google.com/go/v4/auth"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(db.NewDB, cache.NewRedis, infra.NewFirebaseApp, NewData, NewUsersRepo, NewUserSessionRepo)

// Data .
type Data struct {
	query        *query.Query
	redis        *redis.Client
	firebaseAuth *auth.Client
	logger       log.Logger
}

// NewData .
func NewData(db *gorm.DB, redis *redis.Client, firebaseAuth *auth.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Log(log.LevelInfo, "closing the data resources")
	}
	return &Data{
		query:        query.Use(db),
		redis:        redis,
		firebaseAuth: firebaseAuth,
		logger:       logger,
	}, cleanup, nil
}
