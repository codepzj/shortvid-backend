package server

import (
	v1 "shortvid-backend/api/shortvid-service/v1"
	"shortvid-backend/app/shortvid-service/internal/conf"
	"shortvid-backend/app/shortvid-service/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, users *service.UsersService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	// 自定义响应格式
	opts = append(opts, http.ResponseEncoder(ResponseEncoder))
	opts = append(opts, http.ErrorEncoder(ErrorEncoder))

	srv := http.NewServer(opts...)
	v1.RegisterUsersServiceHTTPServer(srv, users)
	return srv
}
