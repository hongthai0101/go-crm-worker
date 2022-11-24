package repositories

import (
	"context"
	"crm-worker-go/datasources"
	"crm-worker-go/entities"
)

type NoteRepository struct {
	BaseRepo *BaseRepository[entities.Note]
}

func NewNoteRepository(ctx context.Context) *NoteRepository {
	return &NoteRepository{
		BaseRepo: &BaseRepository[entities.Note]{
			col: datasources.MongoDatabase.Collection(entities.CollectionNote),
			ctx: ctx,
		},
	}
}
