package service

import (
	"context"
	"log/slog"
	"shortvid-backend/app/user-service/internal/conf"

	"github.com/go-kratos/kratos/v3/log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
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
	auth *auth.Client
}

func NewFirebaseService(c *conf.Firebase) (*FirebaseService, error) {
	auth, err := NewFirebaseApp(c)
	if err != nil {
		return nil, err
	}
	return &FirebaseService{auth: auth}, nil
}

// VertifyIDToken 验证IDToken[验签]
func (s *FirebaseService) VertifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := s.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Error("Verify Firebase ID token failed", slog.Any("error", err))
		return nil, err
	}
	return token, nil
}

// GetUserInfo 获取用户信息
func (s *FirebaseService) GetUserInfo(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := s.auth.GetUser(ctx, uid)
	if err != nil {
		log.Error("Get Firebase User info failed", slog.Any("error", err))
		return nil, err
	}
	return user, nil
}
