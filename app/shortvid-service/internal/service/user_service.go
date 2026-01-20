package service

import (
	"context"
	"errors"
	v1 "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/biz"
	"shortvid-backend/app/shortvid-service/internal/data/model"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type UsersService struct {
	v1.UnimplementedUsersServiceServer

	logger             log.Logger
	uc                 *biz.UsersUsecase
	firebaseService    *FirebaseService
	jwtService         *JwtService
	userSessionService *UserSessionService
	cacheService       *CacheService
}

func NewUsersService(logger log.Logger, uc *biz.UsersUsecase, firebaseService *FirebaseService, jwtService *JwtService, userSessionService *UserSessionService, cacheService *CacheService) *UsersService {
	return &UsersService{
		logger:             logger,
		uc:                 uc,
		firebaseService:    firebaseService,
		jwtService:         jwtService,
		userSessionService: userSessionService,
		cacheService:       cacheService,
	}
}

func (s *UsersService) LoginFirebase(ctx context.Context, req *v1.LoginFirebaseRequest) (*v1.LoginFirebaseResponse, error) {
	// 1. 验证IDToken
	token, err := s.firebaseService.VertifyIDToken(ctx, req.IdToken)
	if err != nil {
		return nil, err
	}

	// 2. 获取用户信息
	providerUID := token.UID
	firebaseUser, err := s.firebaseService.GetUserInfo(ctx, providerUID)
	if err != nil {
		return nil, err
	}
	if firebaseUser == nil {
		return nil, errors.New("firebase user not found")
	}

	// 3. 查找或创建用户
	user, isNew, err := s.uc.FindOrCreateUser(ctx, &biz.User{
		Nickname:    firebaseUser.DisplayName,
		Avatar:      firebaseUser.PhotoURL,
		Email:       firebaseUser.Email,
		ProviderUID: providerUID,
		Provider:    "firebase",
	})
	if err != nil {
		return nil, err
	}

	// 4. 生成sessionID
	sessionID := uuid.NewString()

	// 5. 生成accessToken
	_, err = s.jwtService.GenerateAccessToken(user.UserUID, sessionID)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "Generate access token failed", "error", err)
		return nil, err
	}

	// 6. 生成refreshToken
	_, err = s.jwtService.GenerateRefreshToken(user.UserUID, sessionID)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "Generate refresh token failed", "error", err)
		return nil, err
	}

	// 7. 创建用户会话session
	session := &model.UserSession{
		UserUID:   user.UserUID,
		SessionID: sessionID,
		IP:        "127.0.0.1", // TODO: 从请求中获取真实IP
		UserAgent: "Unknown",   // TODO: 从请求中获取User-Agent
		Platform:  "Unknown",   // TODO: 从请求中获取平台信息
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshTokenExpiration()),
	}
	if err := s.userSessionService.CreateUserSession(ctx, session); err != nil {
		s.logger.Log(log.LevelError, "msg", "Create user session failed", "error", err)
		return nil, err
	}

	// 8. 限制会话数量
	if err := s.userSessionService.LimitUserSession(ctx, user.UserUID); err != nil {
		s.logger.Log(log.LevelError, "msg", "Limit user session failed", "error", err)
		return nil, err
	}

	// 9. 将用户会话缓存到redis中
	expiration := s.jwtService.GetTokenExpiration()
	if err := s.cacheService.SetUserSession(ctx, user.UserUID, sessionID, expiration); err != nil {
		s.logger.Log(log.LevelError, "msg", "Set user session failed", "error", err)
	}

	// 10. 如果用户是新用户，则记录日志
	if isNew {
		s.logger.Log(log.LevelInfo, "msg", "user is new")
	} else {
		s.logger.Log(log.LevelInfo, "msg", "user is found", "user", user)
	}
	return &v1.LoginFirebaseResponse{}, nil
}

func (s *UsersService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	user, err := s.uc.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserResponse{
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Email:       user.Email,
		Provider:    user.Provider,
		ProviderUid: user.ProviderUID,
	}, nil
}
