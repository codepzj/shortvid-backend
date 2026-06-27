package data

import (
	"context"
	"errors"
	"shortvid-backend/app/user-service/internal/biz"
	"shortvid-backend/app/user-service/internal/data/model"

	"gorm.io/gorm"
)

type userSessionRepo struct {
	data *Data
}

func NewUserSessionRepo(data *Data) biz.UserSessionRepo {
	return &userSessionRepo{data: data}
}

func (r *userSessionRepo) CreateUserSession(ctx context.Context, userSession *model.UserSession) error {
	return r.data.db.WithContext(ctx).Create(userSession).Error
}

func (r *userSessionRepo) FindUserSessionByUID(ctx context.Context, UID int) ([]*model.UserSession, error) {
	var sessions []*model.UserSession
	err := r.data.db.WithContext(ctx).Where("uid = ?", UID).Find(&sessions).Error
	return sessions, err
}

func (r *userSessionRepo) FindUserSessionBySessionID(ctx context.Context, sessionID string) (*model.UserSession, error) {
	var session model.UserSession
	err := r.data.db.WithContext(ctx).Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *userSessionRepo) DeleteUserSessionBySessionID(ctx context.Context, sessionID string) error {
	return r.data.db.WithContext(ctx).Where("session_id = ?", sessionID).Delete(&model.UserSession{}).Error
}

func (r *userSessionRepo) DeleteUserSessionByIDs(ctx context.Context, ids []int) error {
	return r.data.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.UserSession{}).Error
}
