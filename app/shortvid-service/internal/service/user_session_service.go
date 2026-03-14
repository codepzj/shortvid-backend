package service

import (
	"context"
	"errors"
	"shortvid-backend/app/shortvid-service/internal/biz"
	"shortvid-backend/app/shortvid-service/internal/conf"
	"shortvid-backend/app/shortvid-service/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

type UserSessionService struct {
	logger          log.Logger
	sessionConf     *conf.Session
	userSessionRepo biz.UserSessionRepo
	cacheService    *CacheService
	jwtService      *JwtService
}

func NewUserSessionService(logger log.Logger, sessionConf *conf.Session, userSessionRepo biz.UserSessionRepo, cacheService *CacheService, jwtService *JwtService) *UserSessionService {
	return &UserSessionService{
		logger:          logger,
		sessionConf:     sessionConf,
		userSessionRepo: userSessionRepo,
		cacheService:    cacheService,
		jwtService:      jwtService,
	}
}

// CreateUserSession 创建用户会话
func (s *UserSessionService) CreateUserSession(ctx context.Context, userSession *model.UserSession) error {
	if err := s.userSessionRepo.CreateUserSession(ctx, userSession); err != nil {
		s.logger.Log(log.LevelError, "msg", "Create user session failed", "error", err)
		return err
	}
	return nil
}

// DeleteUserSession 删除用户会话
func (s *UserSessionService) DeleteUserSession(ctx context.Context, sessionId string) error {
	// 1. 先查数据库
	session, err := s.userSessionRepo.FindUserSessionBySessionID(ctx, sessionId)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "Find user session by session id failed", "error", err)
		return err
	}
	// 2. 如果会话存在, 则删除会话
	if session != nil {
		if err := s.userSessionRepo.DeleteUserSessionBySessionID(ctx, sessionId); err != nil {
			s.logger.Log(log.LevelError, "msg", "Delete user session by session id failed", "error", err)
			return err
		}
	}
	// 3. 删缓存
	s.cacheService.DeleteUserSession(ctx, sessionId)

	return nil
}

// ValidateSession 验证用户会话
func (s *UserSessionService) ValidateSession(ctx context.Context, sessionId string) error {
	// 1. 先查缓存
	sessionInfo, err := s.cacheService.GetUserSessionInfo(ctx, sessionId)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "Get user session info failed", "error", err)
		return err
	}
	// map不为空, 则会话有效
	if len(sessionInfo) > 0 {
		return nil
	}
	// 2. 查询用户会话
	session, err := s.userSessionRepo.FindUserSessionBySessionID(ctx, sessionId)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "Find user session by session id failed", "error", err)
		return err
	}
	// 3. 会话是否过期
	if session == nil {
		return errors.New("session expired")
	}
	// 4. 回填缓存
	expiration := s.jwtService.GetTokenExpiration()
	if err := s.cacheService.SetUserSession(ctx, session.UID, sessionId, expiration); err != nil {
		s.logger.Log(log.LevelError, "msg", "Set user session failed", "error", err)
		return err
	}
	return nil
}

// LimitUserSession 限制用户会话数量
func (s *UserSessionService) LimitUserSession(ctx context.Context, UID int) error {
	// 1. 如果未启用限制，直接返回
	if s.sessionConf == nil || !s.sessionConf.LimitEnabled {
		return nil
	}

	// 2. 获取限制数量
	limitCount := s.sessionConf.LimitCount
	// 如果限制数量小于等于0，则不限制
	if limitCount <= 0 {
		return nil
	}

	// 3. 查询用户会话
	sessions, err := s.userSessionRepo.FindUserSessionByUID(ctx, UID)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "Find user session by user id failed", "error", err)
		return err
	}
	// 4. 会话数量超过限制, 删除最早的会话
	if len(sessions) > int(limitCount) {
		deleteSessions := sessions[limitCount:]
		deleteSessionIds := make([]int, len(deleteSessions))
		for i := range deleteSessions {
			deleteSessionIds[i] = deleteSessions[i].ID
		}
		if err := s.userSessionRepo.DeleteUserSessionByIDs(ctx, deleteSessionIds); err != nil {
			s.logger.Log(log.LevelError, "msg", "Delete user session by ids failed", "error", err)
			return err
		}
		// 删除对应的缓存
		for _, session := range deleteSessions {
			s.cacheService.DeleteUserSession(ctx, session.SessionID)
		}
	}
	return nil
}
