package repositories

import (
	"context"
	"crm-worker-go/datasources"
	"crm-worker-go/entities"
)

type TagRepository struct {
	BaseRepo *BaseRepository[entities.Tag]
}

func NewTagRepository(ctx context.Context) *TagRepository {
	return &TagRepository{
		BaseRepo: &BaseRepository[entities.Tag]{
			col: datasources.MongoDatabase.Collection(entities.CollectionTag),
			ctx: ctx,
		},
	}
}
