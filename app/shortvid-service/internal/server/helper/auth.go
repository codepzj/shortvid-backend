package helper

import (
	"context"
	"errors"
	v1 "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/service"
	"slices"
	"strconv"
	"strings"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// 路由白名单
var authWhiteList = []string{
	v1.OperationUserServiceLoginFirebase,
	v1.OperationUserServiceLoginGithub,
}

func RequireAuthMiddleware(userSvc *service.UserService, jwtSvc *service.JwtService) middleware.Middleware {
	return selector.Server(
		RequireAuth(userSvc, jwtSvc),
	).Match(NewWhiteListMatcher()).Build()
}

func RequireAuth(userSvc *service.UserService, jwtSvc *service.JwtService) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			hr, ok := khttp.RequestFromServerContext(ctx)
			if !ok {
				return nil, errors.New("get http request failed")
			}

			// 开发环境绕过机制
			UID := hr.Header.Get("X-USER-UID")
			if UID != "" {
				UIDInt, err := strconv.Atoi(UID)
				if err != nil {
					return nil, err
				}
				user, err := userSvc.GetUserByUID(ctx, UIDInt)
				if err != nil {
					return nil, err
				}
				ctx = jwtSvc.SetClaimsFromContext(ctx, &service.JwtCustomClaims{
					UID:       user.UID,
					SessionID: "",
				})
				return next(ctx, req)
			}

			authorization := hr.Header.Get("Authorization")

			if authorization == "" {
				return nil, kerrors.Unauthorized("no authorization", "authorization is empty")
			}

			parts := strings.SplitN(authorization, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return nil, kerrors.Unauthorized("authorize failed","invalid authorization format")
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
	return func(ctx context.Context, operation string) bool {
		if slices.Contains(authWhiteList, operation) {
			return false // 不走鉴权中间件
		}
		return true
	}
}
