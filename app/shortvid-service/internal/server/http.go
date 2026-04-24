package server

import (
	v1 "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/conf"
	"shortvid-backend/app/shortvid-service/internal/server/helper"
	"shortvid-backend/app/shortvid-service/internal/service"

	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(cs *conf.Server, cj *conf.Jwt, userSvc *service.UserService, jwtSvc *service.JwtService, fileSvc *service.FileService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.ProtoValidate(), // 校验参数中间件
			helper.RequireAuthMiddleware(userSvc, jwtSvc),
			helper.PublicParamMiddleware(),
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

	srv := http.NewServer(opts...)

	// 用户服务
	v1.RegisterUserServiceHTTPServer(srv, userSvc)
	// 文件服务
	v1.RegisterFileHTTPServer(srv, fileSvc)
	return srv
}
