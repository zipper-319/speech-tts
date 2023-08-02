//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	service1 "speech-tts/internal/cgo/service"
	"speech-tts/internal/conf"
	"speech-tts/internal/data"
	"speech-tts/internal/server"
	"speech-tts/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, service1.ProviderSet, service.ProviderSet, newApp))
}
