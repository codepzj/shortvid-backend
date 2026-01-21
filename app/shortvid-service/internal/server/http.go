package server

import (
	"context"
	v1 "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/conf"
	"shortvid-backend/app/shortvid-service/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	kjwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/golang-jwt/jwt/v5"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(cs *conf.Server, cj *conf.Jwt, users *service.UsersService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if cs.Http.Network != "" {
		opts = append(opts, http.Network(cs.Http.Network))
	}
	if cs.Http.Addr != "" {
		opts = append(opts, http.Address(cs.Http.Addr))
	}
	if cs.Http.Timeout != nil {
		opts = append(opts, http.Timeout(cs.Http.Timeout.AsDuration()))
	}

	// 自定义响应格式
	opts = append(opts, http.ResponseEncoder(ResponseEncoder))

	// 自定义错误格式
	opts = append(opts, http.ErrorEncoder(ErrorEncoder))

	// 中间件
	opts = append(opts, http.Middleware(
		selector.Server(RequireAuth(cj.SecretKey)).Match(NewWhiteListMatcher()).Build(),
	))

	srv := http.NewServer(opts...)
	v1.RegisterUsersServiceHTTPServer(srv, users)
	return srv
}

// RequireAuth 需要认证的接口
func RequireAuth(authKey string) middleware.Middleware {
	return kjwt.Server(func(t *jwt.Token) (any, error) {
		return []byte(authKey), nil
	}, kjwt.WithSigningMethod(jwt.SigningMethodHS256))
}

// NewWhiteListMatcher 白名单跳过鉴权
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList[v1.OperationUsersServiceLoginFirebase] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}
