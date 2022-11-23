// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"crm-worker-go/clients"
	"crm-worker-go/repositories"
	"crm-worker-go/server"
	"crm-worker-go/services"
)

// Injectors from wire.go:

func initServer() *server.Server {
	httpClient := clients.NewHttpClient()
	repository := repositories.NewRepository()
	service := services.NewService(httpClient, repository)
	serverServer := server.NewServer(service)
	return serverServer
}
