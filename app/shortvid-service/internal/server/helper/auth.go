package helper

import (
	"context"
	"errors"
	v1 "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/service"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func RequireAuthMiddleware(userSvc *service.UsersService, jwtSvc *service.JwtService) middleware.Middleware {
	return selector.Server(
		RequireAuth(userSvc, jwtSvc),
	).Match(NewWhiteListMatcher()).Build()
}

// RequireAuth 需要认证的接口
func RequireAuth(userSvc *service.UsersService, jwtSvc *service.JwtService) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			hr, ok := khttp.RequestFromServerContext(ctx)
			if !ok {
				return nil, errors.New("get http request failed")
			}

			// 开发环境绕过机制
			userUID := hr.Header.Get("X-USER-UID")
			if userUID != "" {
				userUIDInt, err := strconv.Atoi(userUID)
				if err != nil {
					return nil, err
				}
				user, err := userSvc.GetUserByUID(ctx, userUIDInt)
				if err != nil {
					return nil, err
				}
				ctx = jwtSvc.SetClaimsFromContext(ctx, &service.JwtCustomClaims{
					UserUID:   user.UserUID,
					SessionID: "",
				})
				return next(ctx, req)
			}

			authorization := hr.Header.Get("Authorization")

			if authorization == "" {
				return nil, errors.New("no authorization")
			}

			parts := strings.SplitN(authorization, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return nil, errors.New("invalid authorization format")
			}

			tokenStr := parts[1]

			claims, err := jwtSvc.ValidateToken(tokenStr)
			if err != nil {
				return nil, err
			}

			ctxWithClaims := jwtSvc.SetClaimsFromContext(ctx, claims)
			return next(ctxWithClaims, req)
		}
	}
}

// NewWhiteListMatcher 白名单跳过鉴权
func NewWhiteListMatcher() selector.MatchFunc {
	// 路由白名单
	whiteList := make(map[string]struct{})
	whiteList[v1.OperationUsersServiceLoginFirebase] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		_, ok := whiteList[operation]
		if ok {
			return false
		}
		return true
	}
}
