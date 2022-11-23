package repositories

import (
	"context"
	"crm-worker-go/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseRepository[T any] struct {
	col *mongo.Collection
	ctx context.Context
}

func (r *BaseRepository[T]) Find(filter interface{}, opts *options.FindOptions) ([]*T, error) {
	var results []*T
	cursor, _ := r.col.Find(r.ctx, filter, opts)

	if err := cursor.All(r.ctx, &results); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return results, nil
}

func (r *BaseRepository[T]) Count(filter interface{}) (int64, error) {
	total, err := r.col.CountDocuments(r.ctx, filter)
	if err != nil {
		utils.Logger.Error(err)
		return 0, err
	}
	return total, nil
}

func (r *BaseRepository[T]) FindOne(filter interface{}, opts *options.FindOneOptions) (*T, error) {
	var item *T

	cursor := r.col.FindOne(r.ctx, filter, opts)
	if cursor.Err() != nil {
		utils.Logger.Error(cursor.Err())
		return nil, cursor.Err()
	}

	if err := cursor.Decode(&item); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	return item, nil
}

func (r *BaseRepository[T]) Create(entity *T) (*T, error) {
	result, err := r.col.InsertOne(r.ctx, entity)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	var item *T
	if err = r.col.FindOne(r.ctx, bson.M{"_id": result.InsertedID}).Decode(&item); err != nil {
		panic(err)
	}
	return item, nil
}

func (r *BaseRepository[T]) UpdateByID(
	Id primitive.ObjectID,
	payload bson.M,
) (*T, error) {
	_, err := r.col.UpdateByID(r.ctx, Id, bson.D{{
		"$set", payload,
	}})
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	var item *T
	if err = r.col.FindOne(r.ctx, bson.M{"_id": Id}).Decode(&item); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return item, nil
}

func (r *BaseRepository[T]) Delete(filter bson.D) (bool, error) {
	_, err := r.col.DeleteOne(r.ctx, filter)
	if err != nil {
		utils.Logger.Error(err)
		return false, err
	}
	return true, nil
}

func (r *BaseRepository[T]) FindById(Id primitive.ObjectID) (*T, error) {
	var item *T

	cursor := r.col.FindOne(r.ctx, bson.M{"_id": Id}, nil)
	if cursor.Err() != nil {
		utils.Logger.Error(cursor.Err())
		return nil, cursor.Err()
	}

	if err := cursor.Decode(&item); err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	return item, nil
}
