//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/bighuangbee/face_search2/app/internal/data"
	"github.com/bighuangbee/face_search2/app/internal/server"
	"github.com/bighuangbee/face_search2/app/internal/service"
	"github.com/bighuangbee/face_search2/pkg/conf"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
)

// wireApp init kratos application.
func wireApp(*conf.Bootstrap, log.Logger, naming_client.INamingClient) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet, newApp, data.ProviderSet))
}
