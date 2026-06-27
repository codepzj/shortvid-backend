//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	"shortvid-backend/app/user-service/internal/biz"
	"shortvid-backend/app/user-service/internal/conf"
	"shortvid-backend/app/user-service/internal/data"
	"shortvid-backend/app/user-service/internal/server"
	"shortvid-backend/app/user-service/internal/service"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Firebase, *conf.Github, *conf.Jwt, *conf.Session, *slog.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
