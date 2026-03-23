package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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

	// 通过code获取accessToken
	accessToken, err := cfg.Exchange(ctx, code)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "get github access token by code failed", "error", err)
		return err
	}
	s.logger.Log(log.LevelDebug, "accessToken", accessToken)

	client := cfg.Client(ctx, accessToken)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		s.logger.Log(log.LevelError, "error", "call github user info api failed")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Log(log.LevelError, "code", resp.StatusCode, "error", "get github user info failed")
		return errors.New("get github user info failed")
	}

	fmt.Println(resp.Body)

	return nil
}
