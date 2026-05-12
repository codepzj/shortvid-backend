package utils

import (
	"context"
	"shortvid-backend/app/shortvid-service/internal/biz"
)

func GetPublicParamFromCtx(ctx context.Context) biz.PublicParam {
	val := ctx.Value("public_param")
	pp, ok := val.(*biz.PublicParam)
	if !ok {
		return biz.PublicParam{}
	}
	return *pp
}
