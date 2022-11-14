package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionNote = "Note"

type Note struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"ID,omitempty"`
	Content             string             `bson:"content" json:"content,omitempty"`
	SaleOpportunitiesId string             `bson:"saleOpportunitiesId" json:"saleOpportunitiesId,omitempty"`
	BaseEntity          `bson:"inline"`
}
