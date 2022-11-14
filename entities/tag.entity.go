package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const CollectionTag = "Tag"

type Tag struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"ID,omitempty"`
	Code      string             `bson:"code" json:"code,omitempty"`
	Name      string             `bson:"name" json:"name,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt,omitempty"`
}
