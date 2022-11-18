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
	storageService := NewStorageService("vmw")
	return &Service{
		ExportService:  NewExportService(client, repository, storageService),
		SaleService:    NewSaleService(NewTopicService(), repository),
		TopicService:   NewTopicService(),
		StorageService: storageService,
	}
}
