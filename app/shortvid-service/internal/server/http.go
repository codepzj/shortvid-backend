package server

import (
	healthV1 "shortvid-backend/api/v1/health"
	uploadV1 "shortvid-backend/api/v1/upload"
	"shortvid-backend/app/shortvid-service/internal/conf"
	"shortvid-backend/app/shortvid-service/internal/service"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/middleware/validate"
	"github.com/go-kratos/kratos/v3/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(cs *conf.Server, healthSvc *service.HealthService, uploadSvc *service.UploadService) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(), // 校验参数中间件
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

	healthV1.RegisterHealthServiceHTTPServer(srv, healthSvc)
	uploadV1.RegisterUploadServiceHTTPServer(srv, uploadSvc)
	return srv
}
