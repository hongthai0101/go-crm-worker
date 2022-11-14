package entities

import (
	time "time"
)

type BaseEntity struct {
	CreatedBy string      `bson:"createdBy" json:"createdBy,omitempty"`
	UpdatedBy string      `bson:"updatedBy" json:"updatedBy,omitempty"`
	CreatedAt time.Time   `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt time.Time   `bson:"updatedAt" json:"updatedAt,omitempty"`
	DeletedAt interface{} `bson:"deletedAt" json:"deletedAt,omitempty"`
}

func CreatingEntity(entity *BaseEntity) {
	entity.UpdatedAt = time.Now()
	entity.CreatedAt = time.Now()
	entity.DeletedAt = nil
}

func UpdatingEntity(entity *BaseEntity) {
	entity.UpdatedAt = time.Now()
}

func DeletingEntity(entity *BaseEntity) {
	entity.DeletedAt = time.Now()
}
