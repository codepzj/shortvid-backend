package service

import (
	"context"
	"errors"
	pb "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/biz"
	"shortvid-backend/app/shortvid-service/internal/data/model"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type UsersService struct {
	pb.UnimplementedUsersServiceServer

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

func (s *UsersService) LoginFirebase(ctx context.Context, req *pb.LoginFirebaseRequest) (*pb.LoginFirebaseResponse, error) {
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

	// 4. 更新登录信息
	if err := s.uc.UpdateLoginInfo(ctx, user.ID); err != nil {
		s.logger.Log(log.LevelError, "msg", "Update login info failed", "error", err)
		return nil, err
	}

	// 5. 生成sessionID
	sessionID := uuid.NewString()

	// 6. 生成accessToken
	accessToken, err := s.jwtService.GenerateAccessToken(user.UserUID, sessionID)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "Generate access token failed", "error", err)
		return nil, err
	}

	// 7. 生成refreshToken
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.UserUID, sessionID)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "Generate refresh token failed", "error", err)
		return nil, err
	}

	// 8. 创建用户会话session
	session := &model.UserSession{
		UserUID:   user.UserUID,
		SessionID: sessionID,
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshTokenExpiration()),
	}
	if err := s.userSessionService.CreateUserSession(ctx, session); err != nil {
		s.logger.Log(log.LevelError, "msg", "Create user session failed", "error", err)
		return nil, err
	}

	// 9. 限制会话数量
	if err := s.userSessionService.LimitUserSession(ctx, user.UserUID); err != nil {
		s.logger.Log(log.LevelError, "msg", "Limit user session failed", "error", err)
		return nil, err
	}

	// 10. 将用户会话缓存到redis中
	expiration := s.jwtService.GetTokenExpiration()
	if err := s.cacheService.SetUserSession(ctx, user.UserUID, sessionID, expiration); err != nil {
		s.logger.Log(log.LevelError, "msg", "Set user session failed", "error", err)
	}

	// 11. 如果用户是新用户，则记录日志
	if isNew {
		s.logger.Log(log.LevelInfo, "msg", "user is new")
	} else {
		s.logger.Log(log.LevelInfo, "msg", "user already esist", "user", user)
	}

	return &pb.LoginFirebaseResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &pb.UserProfile{
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		},
	}, nil
}

// GetUserProfile 根据userUid查询用户信息
func (s *UsersService) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	user, err := s.uc.GetUserByUserUID(ctx, int(req.UserUid))
	if err != nil {
		return nil, err
	}
	return &pb.GetUserProfileResponse{
		UserInfo: &pb.UserProfile{
			UserUid:  int32(user.UserUID),
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		},
	}, nil
}

func (s *UsersService) UserInfo(ctx context.Context, req *emptypb.Empty) (*pb.UserInfoResponse, error) {
	claims, err := s.jwtService.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	user, err := s.uc.GetUserByUserUID(ctx, claims.UserUID)
	if err != nil {
		return nil, err
	}
	return &pb.UserInfoResponse{
		UserInfo: &pb.UserProfile{
			Id:       int32(user.ID),
			UserUid:  int32(user.UserUID),
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		},
	}, nil
}

// GetUserByUID 根据UID查询用户
func (s *UsersService) GetUserByUserUID(ctx context.Context, userUID int) (*biz.UserProfile, error) {
	user, err := s.uc.GetUserByUserUID(ctx, userUID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
