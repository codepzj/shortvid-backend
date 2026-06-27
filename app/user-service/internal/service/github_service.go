package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"shortvid-backend/app/user-service/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v3/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubUserInfo struct {
	AvatarURL         string    `json:"avatar_url"`
	Bio               string    `json:"bio"`
	Company           string    `json:"company"`
	CreatedAt         time.Time `json:"created_at"`
	Email             *string   `json:"email"`
	EventsURL         string    `json:"events_url"`
	Followers         int       `json:"followers"`
	FollowersURL      string    `json:"followers_url"`
	Following         int       `json:"following"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	GravatarID        string    `json:"gravatar_id"`
	Hireable          *bool     `json:"hireable"`
	HTMLURL           string    `json:"html_url"`
	ID                int       `json:"id"`
	Location          *string   `json:"location"`
	Login             string    `json:"login"`
	Name              string    `json:"name"`
	NodeID            string    `json:"node_id"`
	OrganizationsURL  string    `json:"organizations_url"`
	PublicGists       int       `json:"public_gists"`
	PublicRepos       int       `json:"public_repos"`
	ReceivedEventsURL string    `json:"received_events_url"`
	ReposURL          string    `json:"repos_url"`
	SiteAdmin         bool      `json:"site_admin"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	TwitterUsername   *string   `json:"twitter_username"`
	Type              string    `json:"type"`
	UpdatedAt         time.Time `json:"updated_at"`
	URL               string    `json:"url"`
}

type GithubService struct {
	ghCfg *conf.Github
}

func NewGithubService(ghCfg *conf.Github) *GithubService {
	return &GithubService{
		ghCfg: ghCfg,
	}
}

// GetGithubUserInfo 获取Github用户信息
func (s *GithubService) GetGithubUserInfo(ctx context.Context, code string) (*GithubUserInfo, error) {
	cfg := oauth2.Config{
		ClientID:     s.ghCfg.ClientId,
		ClientSecret: s.ghCfg.ClientSecret,
		Scopes:       []string{"user"},
		Endpoint:     github.Endpoint,
	}

	// 通过code获取accessToken
	accessToken, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Error("get github access token by code failed", slog.Any("error", err))
		return nil, err
	}

	// 通过accessToken获取用户信息
	client := cfg.Client(ctx, accessToken)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo GithubUserInfo
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	log.Debug("userInfo", slog.Any("userInfo", userInfo))
	return &userInfo, nil
}
