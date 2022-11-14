package repositories

import (
	"context"
	"crm-worker-go/config"
	"crm-worker-go/datasources"
	"crm-worker-go/entities"
)

type LeadRepository struct {
	BaseRepo *BaseRepository[entities.Lead]
}

func NewLeadRepository(ctx context.Context) *LeadRepository {
	return &LeadRepository{
		BaseRepo: &BaseRepository[entities.Lead]{
			col: datasources.MongoClient.Database(config.GetConfigDB().Name).Collection(entities.CollectionLead),
			ctx: ctx,
		},
	}
}
