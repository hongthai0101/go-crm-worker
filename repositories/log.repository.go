package repositories

import (
	"context"
	"crm-worker-go/config"
	"crm-worker-go/datasources"
	"crm-worker-go/entities"
)

type LogRepository struct {
	BaseRepo *BaseRepository[entities.Log]
}

func NewLogRepository(ctx context.Context) *LogRepository {
	return &LogRepository{
		BaseRepo: &BaseRepository[entities.Log]{
			col: datasources.MongoClient.Database(config.GetConfigDB().Name).Collection(entities.CollectionLog),
			ctx: ctx,
		},
	}
}
