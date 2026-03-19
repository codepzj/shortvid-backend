package service

import (
	"context"
	"errors"
	"fmt"
	pb "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/biz"
	"shortvid-backend/app/shortvid-service/internal/data/model"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer

	logger             log.Logger
	uc                 *biz.UsersUsecase
	firebaseService    *FirebaseService
	githubService      *GithubService
	jwtService         *JwtService
	userSessionService *UserSessionService
	cacheService       *CacheService
}

func NewUserService(logger log.Logger, uc *biz.UsersUsecase, firebaseService *FirebaseService, githubService *GithubService, jwtService *JwtService, userSessionService *UserSessionService, cacheService *CacheService) *UserService {
	return &UserService{
		logger:             logger,
		uc:                 uc,
		firebaseService:    firebaseService,
		githubService:      githubService,
		jwtService:         jwtService,
		userSessionService: userSessionService,
		cacheService:       cacheService,
	}
}

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
	user, isNew, err := s.uc.FirebaseFindOrCreateUser(ctx, &biz.UserDTO{
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
		s.logger.Log(log.LevelError, "msg", "update login info failed", "error", err)
		return nil, err
	}

	// 5. 生成sessionID
	sessionID := uuid.NewString()

	// 6. 生成accessToken
	accessToken, err := s.jwtService.GenerateAccessToken(user.UID, sessionID)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "generate access token failed", "error", err)
		return nil, err
	}

	// 7. 生成refreshToken
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.UID, sessionID)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "generate refresh token failed", "error", err)
		return nil, err
	}

	// 8. 创建用户会话session
	session := &model.UserSession{
		UID:       user.UID,
		SessionID: sessionID,
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshTokenExpiration()),
	}
	if err := s.userSessionService.CreateUserSession(ctx, session); err != nil {
		s.logger.Log(log.LevelError, "msg", "create user session failed", "error", err)
		return nil, err
	}

	// 9. 限制会话数量
	if err := s.userSessionService.LimitUserSession(ctx, user.UID); err != nil {
		s.logger.Log(log.LevelError, "msg", "Limit user session failed", "error", err)
		return nil, err
	}

	// 10. 将用户会话缓存到redis中
	expiration := s.jwtService.GetTokenExpiration()
	if err := s.cacheService.SetUserSession(ctx, user.UID, sessionID, expiration); err != nil {
		s.logger.Log(log.LevelError, "msg", "set user session failed", "error", err)
	}

	// 11. 如果用户是新用户，则记录日志
	if isNew {
		s.logger.Log(log.LevelInfo, "msg", "user is new")
	} else {
		s.logger.Log(log.LevelInfo, "msg", "user already esist", "user", user)
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

func (s *UserService) LoginGithub(ctx context.Context, req *pb.GithubLoginRequest) (*pb.GithubLoginResponse, error) {
	code := req.Code
	fmt.Println("code", code)

	s.githubService.GetGithubUserInfo(ctx, code)

	return &pb.GithubLoginResponse{
		AccessToken:  "432",
		RefreshToken: "5ewr",
		User:         &pb.UserProfile{},
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
