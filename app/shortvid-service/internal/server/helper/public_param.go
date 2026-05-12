package helper

import (
	"context"
	"errors"
	"shortvid-backend/app/shortvid-service/internal/biz"

	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func PublicParamMiddleware() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			hr, ok := khttp.RequestFromServerContext(ctx)
			if !ok {
				return nil, errors.New("get http request failed")
			}

			userAgent := hr.Header.Get("User-Agent")
			ip := hr.Header.Get("ip")

			pp := &biz.PublicParam{
				UserAgent: userAgent,
				IP:        ip,
			}
			ctx = context.WithValue(ctx, "public_param", pp) // 写入ctx
			return next(ctx, req)
		}
	}
}
