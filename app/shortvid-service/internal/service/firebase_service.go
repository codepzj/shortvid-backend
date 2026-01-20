package service

import (
	"context"
	"shortvid-backend/app/shortvid-service/internal/conf"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/api/option"
)

func NewFirebaseApp(c *conf.Firebase) (*auth.Client, error) {
	opt := option.WithCredentialsJSON([]byte(c.CredentialsJson))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}
	return client, nil
}

type FirebaseService struct {
	logger log.Logger
	auth   *auth.Client
}

func NewFirebaseService(logger log.Logger, c *conf.Firebase) (*FirebaseService, error) {
	auth, err := NewFirebaseApp(c)
	if err != nil {
		logger.Log(log.LevelError, "msg", "New Firebase App failed", "error", err)
		return nil, err
	}
	return &FirebaseService{logger: logger, auth: auth}, nil
}

// VertifyIDToken 验证IDToken[验签]
func (s *FirebaseService) VertifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := s.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "Verify Firebase ID token failed", "error", err)
		return nil, err
	}
	return token, nil
}

// GetUserInfo 获取用户信息
func (s *FirebaseService) GetUserInfo(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := s.auth.GetUser(ctx, uid)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "Get Firebase User info failed", "error", err)
		return nil, err
	}
	return user, nil
}
