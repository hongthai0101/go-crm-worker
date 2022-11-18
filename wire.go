//go:build wireinject
// +build wireinject

package main

import (
	"crm-worker-go/clients"
	"crm-worker-go/repositories"
	"crm-worker-go/server"
	"crm-worker-go/services"
	"github.com/google/wire"
)

func initServer() *server.Server {
	wire.Build(
		repositories.ProviderRepositorySet,
		clients.ProviderHttpClientSet,
		services.ProviderServiceSet,
		server.ProviderServerSet,
	)
	return &server.Server{}
}
