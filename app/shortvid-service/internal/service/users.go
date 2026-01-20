package service

import (
	"context"
	v1 "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type UsersService struct {
	v1.UnimplementedUsersServiceServer

	logger          log.Logger
	uc              *biz.UsersUsecase
	firebaseService *FirebaseService
}

func NewUsersService(logger log.Logger, uc *biz.UsersUsecase, firebaseService *FirebaseService) *UsersService {
	return &UsersService{logger: logger, uc: uc, firebaseService: firebaseService}
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


	// 3. 查找或创建用户
	user, isNew, err := s.uc.FindOrCreateUser(ctx, &biz.User{
		UserUID:     generateUserUID(), // 生成用户唯一ID
		Nickname:    firebaseUser.DisplayName,
		Avatar:      firebaseUser.PhotoURL,
		Email:       firebaseUser.Email,
		ProviderUID: providerUID,
		Provider:    "firebase",
	})
	if err != nil {
		return nil, err
	}

	// 4. 如果用户是新用户，则记录日志
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

func generateUserUID() string {
	return uuid.NewString()
}
