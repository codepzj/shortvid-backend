package service

import (
	"context"

	"firebase.google.com/go/v4/auth"
	firebase "firebase.google.com/go/v4/auth"
	"github.com/go-kratos/kratos/v2/log"
)

type FirebaseService struct {
	logger log.Logger
	auth   *auth.Client
}

func NewFirebaseService(logger log.Logger, auth *auth.Client) *FirebaseService {
	return &FirebaseService{auth: auth, logger: logger}
}

// VertifyIDToken 验证IDToken[验签]
func (s *FirebaseService) VertifyIDToken(ctx context.Context, idToken string) (*firebase.Token, error) {
	token, err := s.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "Verify Firebase ID token failed", "error", err)
		return nil, err
	}
	return token, nil
}

// GetUserInfo 获取用户信息
func (s *FirebaseService) GetUserInfo(ctx context.Context, uid string) (*firebase.UserRecord, error) {
	user, err := s.auth.GetUser(ctx, uid)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "Get Firebase User info failed", "error", err)
		return nil, err
	}
	return user, nil
}
