package repositories

import (
	"context"
	"crm-worker-go/config"
	"crm-worker-go/datasources"
	"crm-worker-go/entities"
)

type TagRepository struct {
	BaseRepo *BaseRepository[entities.Tag]
}

func NewTagRepository(ctx context.Context) *TagRepository {
	return &TagRepository{
		BaseRepo: &BaseRepository[entities.Tag]{
			col: datasources.MongoClient.Database(config.GetConfig().DB.Name).Collection(entities.CollectionTag),
			ctx: ctx,
		},
	}
}
