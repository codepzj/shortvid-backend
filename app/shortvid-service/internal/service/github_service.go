package service

import (
	"context"
	"shortvid-backend/app/shortvid-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubService struct {
	logger log.Logger
	ghCfg  *conf.Github
}

func NewGithubService(logger log.Logger, ghCfg *conf.Github) *GithubService {
	return &GithubService{
		logger: logger,
		ghCfg:  ghCfg,
	}
}

// GetGithubUserInfo 获取Github用户信息
func (s *GithubService) GetGithubUserInfo(ctx context.Context, code string) error {
	cfg := oauth2.Config{
		ClientID:     s.ghCfg.ClientId,
		ClientSecret: s.ghCfg.ClientSecret,
		Scopes:       []string{"user"},
		Endpoint:     github.Endpoint,
	}
	accessToken, err := cfg.Exchange(ctx, code)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "get github access token by code failed", "error", err)
		return err
	}
	s.logger.Log(log.LevelDebug, "accessToken", accessToken)

	return nil
}
