package repositories

import (
	"context"
	"crm-worker-go/datasources"
	"crm-worker-go/entities"
)

type LeadRepository struct {
	BaseRepo *BaseRepository[entities.Lead]
}

func NewLeadRepository(ctx context.Context) *LeadRepository {
	return &LeadRepository{
		BaseRepo: &BaseRepository[entities.Lead]{
			col: datasources.MongoDatabase.Collection(entities.CollectionLead),
			ctx: ctx,
		},
	}
}
