package service

import (
	"context"
	"shortvid-backend/app/shortvid-service/internal/data"

	firebase "firebase.google.com/go/v4/auth"
	"github.com/go-kratos/kratos/v2/log"
)

type FirebaseService struct {
	logger log.Logger
	data   *data.Data
}

func NewFirebaseService(logger log.Logger, data *data.Data) *FirebaseService {
	return &FirebaseService{data: data, logger: logger}
}

// VertifyIDToken 验证IDToken[验签]
func (s *FirebaseService) VertifyIDToken(ctx context.Context, idToken string) (*firebase.Token, error) {
	auth := s.data.GetFirebaseAuth()
	token, err := auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "Verify Firebase ID token failed", "error", err)
		return nil, err
	}
	return token, nil
}

// GetUserInfo 获取用户信息
func (s *FirebaseService) GetUserInfo(ctx context.Context, uid string) (*firebase.UserRecord, error) {
	auth := s.data.GetFirebaseAuth()
	user, err := auth.GetUser(ctx, uid)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "Get Firebase User info failed", "error", err)
		return nil, err
	}
	return user, nil
}
