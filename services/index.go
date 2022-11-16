package services

import (
	"crm-worker-go/clients"
	"crm-worker-go/repositories"
)

type Service struct {
	ExportService  *ExportService
	SaleService    *SaleService
	TopicService   *TopicService
	StorageService StorageService
}

func NewService(
	client *clients.HttpClient,
	repository *repositories.Repository,
) *Service {
	return &Service{
		ExportService:  NewExportService(client, repository),
		SaleService:    NewSaleService(NewTopicService(), repository),
		TopicService:   NewTopicService(),
		StorageService: NewStorageService("vmw"),
	}
}
