package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"shortvid-backend/app/user-service/internal/data/infra/cache"
	"time"

	"github.com/go-kratos/kratos/v3/log"
	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	cache *redis.Client
}

func NewCacheService(cache *redis.Client) *CacheService {
	return &CacheService{cache: cache}
}

// SetUserSession 设置用户会话
// sessionID:UID映射[string] & sessionID:session映射[hash]
func (s *CacheService) SetUserSession(ctx context.Context, UID int, sessionID string, expiration time.Duration) error {
	// 1. 设置sessionID和UID的映射关系
	sessionUserKey := cache.GetSessionUserKey(sessionID)
	err := s.cache.Set(ctx, sessionUserKey, UID, expiration).Err()
	if err != nil {
		log.Error("Set session user failed", slog.Any("error", err))
		return err
	}

	// 2. 设置sessionID和session的映射关系
	userSessionKey := cache.GetUserSessionKey(sessionID)
	pipe := s.cache.Pipeline()
	sessionData := map[string]any{
		"uid":        UID,
		"session_id": sessionID,
		"created_at": time.Now().Unix(),
	}
	pipe.HSet(ctx, userSessionKey, sessionData)
	pipe.Expire(ctx, userSessionKey, expiration)
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Error("Set user session failed", slog.Any("error", err))
		return err
	}
	return nil
}

// DeleteUserSession 删除用户会话
func (s *CacheService) DeleteUserSession(ctx context.Context, sessionID string) {
	sessionUserKey := cache.GetSessionUserKey(sessionID)
	err := s.cache.Del(ctx, sessionUserKey).Err()
	if err != nil {
		log.Error("Delete session user key failed", slog.Any("error", err), slog.String("sessionID", sessionID))
	}
	userSessionKey := cache.GetUserSessionKey(sessionID)
	err = s.cache.Del(ctx, userSessionKey).Err()
	if err != nil {
		log.Error("Delete user session key failed", slog.Any("error", err), slog.String("sessionID", sessionID))
	}
}

// GetUserSessionInfo 获取用户会话信息
func (s *CacheService) GetUserSessionInfo(ctx context.Context, sessionID string) (map[string]string, error) {
	userSessionKey := cache.GetUserSessionKey(sessionID)
	return s.cache.HGetAll(ctx, userSessionKey).Result()
}

// SetUserInfo 设置用户信息
func (s *CacheService) SetUserInfo(ctx context.Context, userUID int, userInfo map[string]any, expiration time.Duration) error {
	userInfoKey := cache.GetUserInfoKey(userUID)

	// 序列化userInfo
	userInfoJson, err := json.Marshal(userInfo)
	if err != nil {
		log.Error("Marshal user info failed", slog.Any("error", err))
		return err
	}

	// 设置userInfo
	err = s.cache.Set(ctx, userInfoKey, userInfoJson, expiration).Err()
	if err != nil {
		log.Error("Set user info failed", slog.Any("error", err))
		return err
	}
	return nil
}
