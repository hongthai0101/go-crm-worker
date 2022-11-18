package server

import (
	"crm-worker-go/clients"
	"crm-worker-go/repositories"
	"crm-worker-go/services"
	"crm-worker-go/subscriptions"
	"crm-worker-go/types"
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

func (s *Server) Run() {
	println("Start Server")

	payload := types.PayloadBorrowDisbursed{
		LoanAmount:     6000000,
		ModifiedAmount: 2000000,
		ContractCode:   "CLLQ01129691",
	}

	s.Service.SaleService.Disbursed(payload)

	subscription := subscriptions.NewSubscription(s.Service)
	subscription.Boot()
}

func (s *Server) Close() error {
	println("Close Server")
	return nil
}
