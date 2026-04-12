package helper

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func CommonParamMiddleware() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			hr, ok := khttp.RequestFromServerContext(ctx)
			if !ok {
				return nil, errors.New("get http request failed")
			}

			fmt.Println("hr", hr)
			return next(ctx, req)
		}
	}
}
