package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	pb "shortvid-backend/api/user-service/v1"
	"shortvid-backend/app/user-service/internal/biz"
	"shortvid-backend/app/user-service/internal/data/model"
	"shortvid-backend/app/user-service/pkg/utils"
	"time"

	"github.com/go-kratos/kratos/v3/log"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer

	uc                 *biz.UsersUsecase
	firebaseService    *FirebaseService
	githubService      *GithubService
	jwtService         *JwtService
	userSessionService *UserSessionService
	cacheService       *CacheService
}

func NewUserService(uc *biz.UsersUsecase, firebaseService *FirebaseService, githubService *GithubService, jwtService *JwtService, userSessionService *UserSessionService, cacheService *CacheService) *UserService {
	return &UserService{
		uc:                 uc,
		firebaseService:    firebaseService,
		githubService:      githubService,
		jwtService:         jwtService,
		userSessionService: userSessionService,
		cacheService:       cacheService,
	}
}

// LoginFirebase 登录Firebase
func (s *UserService) LoginFirebase(ctx context.Context, req *pb.FirebaseLoginRequest) (*pb.FirebaseLoginResponse, error) {
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
	user, isNew, err := s.uc.FindOrCreateUser(ctx, &biz.UserDTO{
		Nickname:    firebaseUser.DisplayName,
		Avatar:      firebaseUser.PhotoURL,
		Email:       firebaseUser.Email,
		Provider:    "firebase",
		ProviderUID: providerUID,
	})
	if err != nil {
		return nil, err
	}

	// 4. 更新登录信息
	if err := s.uc.UpdateLoginInfo(ctx, user.UID); err != nil {
		log.Error("update login info failed", slog.Any("error", err))
		return nil, err
	}

	// 5. 生成sessionID
	sessionID := uuid.NewString()

	// 6. 生成accessToken
	accessToken, err := s.jwtService.GenerateAccessToken(user.UID, sessionID)
	if err != nil {
		log.Error("generate access token failed", slog.Any("error", err))
		return nil, err
	}

	// 7. 生成refreshToken
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.UID, sessionID)
	if err != nil {
		log.Error("generate refresh token failed", slog.Any("error", err))
		return nil, err
	}

	// 8. 创建用户会话session
	session := &model.UserSession{
		UID:       user.UID,
		SessionID: sessionID,
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshTokenExpiration()),
	}
	if err := s.userSessionService.CreateUserSession(ctx, session); err != nil {
		log.Error("create user session failed", slog.Any("error", err))
		return nil, err
	}

	// 9. 限制会话数量
	if err := s.userSessionService.LimitUserSession(ctx, user.UID); err != nil {
		log.Error("Limit user session failed", slog.Any("error", err))
		return nil, err
	}

	// 10. 将用户会话缓存到redis中
	expiration := s.jwtService.GetTokenExpiration()
	if err := s.cacheService.SetUserSession(ctx, user.UID, sessionID, expiration); err != nil {
		log.Error("set user session failed", slog.Any("error", err))
	}

	// 11. 如果用户是新用户，则记录日志
	if isNew {
		log.Info("user is new", slog.Any("user", user))
	} else {
		log.Info("user already esist", slog.Any("user", user))
	}

	return &pb.FirebaseLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &pb.UserProfile{
			Id:          int32(user.UID),
			Uid:         int32(user.UID),
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			Email:       user.Email,
			Provider:    user.Provider,
			ProviderUid: user.ProviderUID,
		},
	}, nil
}

// LoginGithub 登录Github
func (s *UserService) LoginGithub(ctx context.Context, req *pb.GithubLoginRequest) (*pb.GithubLoginResponse, error) {
	// 1. 获取Github用户信息
	userInfo, err := s.githubService.GetGithubUserInfo(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	// 2. 查找或创建用户
	user, isNew, err := s.uc.FindOrCreateUser(ctx, &biz.UserDTO{
		Nickname:    userInfo.Name,
		Avatar:      userInfo.AvatarURL,
		Email:       *userInfo.Email,
		Provider:    "github",
		ProviderUID: fmt.Sprintf("%d", userInfo.ID),
	})
	if err != nil {
		return nil, err
	}

	// 3. 更新登录信息
	_ = utils.GetPublicParamFromCtx(ctx)
	if err := s.uc.UpdateLoginInfo(ctx, user.UID); err != nil {
		log.Error("update login info failed", slog.Any("error", err))
		return nil, err
	}

	// 4. 生成sessionID
	sessionID := uuid.NewString()

	// 5. 生成accessToken
	accessToken, err := s.jwtService.GenerateAccessToken(user.UID, sessionID)
	if err != nil {
		log.Error("generate access token failed", slog.Any("error", err))
		return nil, err
	}

	// 6. 生成refreshToken
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.UID, sessionID)
	if err != nil {
		log.Error("generate refresh token failed", slog.Any("error", err))
		return nil, err
	}

	// 7. 创建用户会话session
	session := &model.UserSession{
		UID:       user.UID,
		SessionID: sessionID,
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshTokenExpiration()),
	}
	if err := s.userSessionService.CreateUserSession(ctx, session); err != nil {
		log.Error("create user session failed", slog.Any("error", err))
		return nil, err
	}

	// 8. 限制会话数量
	if err := s.userSessionService.LimitUserSession(ctx, user.UID); err != nil {
		log.Error("Limit user session failed", slog.Any("error", err))
		return nil, err
	}

	// 9. 将用户会话缓存到redis中
	expiration := s.jwtService.GetTokenExpiration()
	if err := s.cacheService.SetUserSession(ctx, user.UID, sessionID, expiration); err != nil {
		log.Error("set user session failed", slog.Any("error", err))
	}

	// 10. 如果用户是新用户，则记录日志
	if isNew {
		log.Info("user is new", slog.Any("user", user))
	} else {
		log.Info("user already esist", slog.Any("user", user))
	}

	return &pb.GithubLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &pb.UserProfile{
			Id:          int32(user.ID),
			Uid:         int32(user.UID),
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			Email:       user.Email,
			Provider:    user.Provider,
			ProviderUid: user.ProviderUID,
		},
	}, nil
}

// GetUserProfile 根据uid查询用户信息
func (s *UserService) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	user, err := s.uc.GetUserByUID(ctx, int(req.Uid))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetUserProfileResponse{
		User: &pb.UserProfile{
			Uid:      int32(user.UID),
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		},
	}, nil
}

func (s *UserService) UserInfo(ctx context.Context, req *emptypb.Empty) (*pb.UserInfoResponse, error) {
	claims, err := s.jwtService.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	user, err := s.uc.GetUserByUID(ctx, claims.UID)
	if err != nil {
		return nil, err
	}
	return &pb.UserInfoResponse{
		User: &pb.UserProfile{
			Id:       int32(user.ID),
			Uid:      int32(user.UID),
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		},
	}, nil
}

// GetUserByUID 根据UID查询用户
func (s *UserService) GetUserByUID(ctx context.Context, uid int) (*biz.UserProfileVO, error) {
	user, err := s.uc.GetUserByUID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}
