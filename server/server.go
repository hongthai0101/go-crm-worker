package server

import (
	"crm-worker-go/services"
	"crm-worker-go/subscriptions"
)

type Server struct {
	service *services.Service
}

func NewServer(
	service *services.Service,
) *Server {
	return &Server{
		service,
	}
}

func (s *Server) Run() {
	println("Start Server")

	//payload := types.MessageOrderDisbursed{
	//	LoanAmount:     6000000,
	//	ModifiedAmount: 2000000,
	//	ContractCode:   "CLLQ01129691",
	//}
	//
	//s.Service.SaleService.OrderDisbursed(payload)

	subscription := subscriptions.NewSubscription(s.service)
	subscription.Boot()
}

func (s *Server) Close() error {
	println("Close Server")
	return nil
}
