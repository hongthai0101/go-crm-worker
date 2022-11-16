package server

import (
	"crm-worker-go/clients"
	"crm-worker-go/repositories"
	"crm-worker-go/services"
	"crm-worker-go/subscriptions"
)

type Server struct {
	Repo       *repositories.Repository
	HttpClient *clients.HttpClient
	Service    *services.Service
}

func NewServer(
	repo *repositories.Repository,
	httpClient *clients.HttpClient,
	service *services.Service,
) *Server {
	return &Server{
		repo,
		httpClient,
		service,
	}
}

func (s *Server) Run() error {

	subscription := subscriptions.NewSubscription(s.Service)
	subscription.Boot()

	println("Start Server")
	return nil
}

func (s *Server) Close() error {
	println("Close Server")
	return nil
}
