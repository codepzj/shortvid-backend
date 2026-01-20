package data

import (
	"context"
	"errors"
	"shortvid-backend/app/shortvid-service/internal/biz"
	"shortvid-backend/app/shortvid-service/internal/data/model"

	"gorm.io/gorm"
)

type userSessionRepo struct {
	data *Data
}

func NewUserSessionRepo(data *Data) biz.UserSessionRepo {
	return &userSessionRepo{data: data}
}

func (r *userSessionRepo) CreateUserSession(ctx context.Context, userSession *model.UserSession) error {
	return r.data.query.UserSession.WithContext(ctx).Create(userSession)
}

func (r *userSessionRepo) FindUserSessionByUserUID(ctx context.Context, userUID int32) ([]*model.UserSession, error) {
	return r.data.query.UserSession.WithContext(ctx).Where(r.data.query.UserSession.UserUID.Eq(userUID)).Find()
}

func (r *userSessionRepo) FindUserSessionBySessionID(ctx context.Context, sessionID string) (*model.UserSession, error) {
	session, err := r.data.query.UserSession.WithContext(ctx).Where(r.data.query.UserSession.SessionID.Eq(sessionID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return session, nil
}

func (r *userSessionRepo) DeleteUserSessionBySessionID(ctx context.Context, sessionID string) error {
	_, err := r.data.query.UserSession.WithContext(ctx).Where(r.data.query.UserSession.SessionID.Eq(sessionID)).Delete()
	return err
}

func (r *userSessionRepo) DeleteUserSessionByIDs(ctx context.Context, ids []int64) error {
	_, err := r.data.query.UserSession.WithContext(ctx).Where(r.data.query.UserSession.ID.In(ids...)).Delete()
	return err
}
