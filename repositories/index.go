package repositories

import "context"

type Repository struct {
	LeadRepo *LeadRepository
	SaleRepo *SaleOpportunityRepository
	LogRepo  *LogRepository
	NoteRepo *NoteRepository
	TagRepo  *TagRepository
}

func NewRepository() *Repository {
	ctx := context.Background()
	return &Repository{
		LeadRepo: NewLeadRepository(ctx),
		SaleRepo: NewSaleOpportunityRepository(ctx),
		LogRepo:  NewLogRepository(ctx),
		NoteRepo: NewNoteRepository(ctx),
		TagRepo:  NewTagRepository(ctx),
	}
}
