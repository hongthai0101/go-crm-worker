package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const CollectionLog = "Log"

type BeforeAttributes struct {
}

type AfterAttributes struct {
}

type Log struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"ID,omitempty"`
	BeforeAttributes    BeforeAttributes   `bson:"beforeAttributes" json:"beforeAttributes"`
	AfterAttributes     AfterAttributes    `bson:"afterAttributes" json:"afterAttributes"`
	SaleOpportunitiesId string             `bson:"saleOpportunitiesId" json:"saleOpportunitiesId,omitempty"`
	CreatedBy           string             `bson:"createdBy" json:"createdBy,omitempty"`
	CreatedAt           time.Time          `bson:"createdAt" json:"createdAt,omitempty"`
}
