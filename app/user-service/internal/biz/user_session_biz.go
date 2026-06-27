package biz

import (
	"context"
	"shortvid-backend/app/user-service/internal/data/model"
)

type UserSessionUsecase struct {
	repo UserSessionRepo
}

func NewUserSessionUsecase(repo UserSessionRepo) *UserSessionUsecase {
	return &UserSessionUsecase{repo: repo}
}

type UserSessionRepo interface {
	CreateUserSession(ctx context.Context, userSession *model.UserSession) error
	FindUserSessionByUID(ctx context.Context, UID int) ([]*model.UserSession, error)
	FindUserSessionBySessionID(ctx context.Context, sessionID string) (*model.UserSession, error)
	DeleteUserSessionBySessionID(ctx context.Context, sessionID string) error
	DeleteUserSessionByIDs(ctx context.Context, ids []int) error
}
