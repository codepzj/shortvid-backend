package service

import (
	"context"
	"io"
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
		RedirectURL:  "http://localhost:5173/login/github/callback",
		Scopes:       []string{"user"},
		Endpoint:     github.Endpoint,
	}

	s.logger.Log(log.LevelDebug, "cfg", cfg, "code", code)

	// 通过code获取accessToken
	accessToken, err := cfg.Exchange(ctx, code)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "get github access token by code failed", "error", err)
		return err
	}
	s.logger.Log(log.LevelDebug, "accessToken", accessToken)

	userURL := "https://api.github.com/user"
	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken.AccessToken)
	req.Header.Add("User-Agent", "Go OAuth App") // GitHub API requires a User-Agent header
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	s.logger.Log(log.LevelDebug, "body", string(body))

	return nil
}
