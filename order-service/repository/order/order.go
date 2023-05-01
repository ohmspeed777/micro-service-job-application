package order

import (
	"app/internal/core/domain"
	"app/internal/core/ports"
	"app/repository"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName = "orders"
)

type Repo struct {
	repository.Repo
}

func NewRepo(db *mongo.Database) ports.IOrderRepo {
	return &Repo{
		Repo: repository.Repo{
			Collection: db.Collection(collectionName),
		},
	}
}

func (r *Repo) FindOneByID(id string) (*domain.Order, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	entity := &domain.Order{}
	err = r.Collection.FindOne(ctx, primitive.M{"_id": _id}).Decode(entity)

	return entity, err
}

func (r *Repo) AggregateOneByID(id string) (*domain.OrderLookedUp, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	entity := []*domain.OrderLookedUp{}

	filter := primitive.M{"_id": _id}

	match := primitive.M{
		"$match": filter,
	}

	lookup := primitive.M{
		"$lookup": primitive.M{
			"from":         "products",
			"localField":   "items.product_id",
			"foreignField": "_id",
			"as":           "products_joined",
		},
	}

	pipe := []primitive.M{match, lookup}

	cursor, err := r.Collection.Aggregate(ctx, pipe)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &entity)
	if err != nil {
		return nil, err
	}

	if len(entity) <= 0 {
		return nil, mongo.ErrNoDocuments
	}

	return entity[0], err
}

func (r *Repo) AggregateAllByUser(userId string, q domain.Query) ([]*domain.OrderLookedUp, int64, error) {
	_id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	entity := []*domain.OrderLookedUp{}
	counter := []domain.Counter{}

	filter := primitive.M{"user": _id}

	match := primitive.M{
		"$match": filter,
	}

	lookup := primitive.M{
		"$lookup": primitive.M{
			"from":         "products",
			"localField":   "items.product_id",
			"foreignField": "_id",
			"as":           "products_joined",
		},
	}

	skip := primitive.M{
		"$skip": q.GetSkip(),
	}

	limit := primitive.M{
		"$limit": q.GetLimit(),
	}

	pipe := []primitive.M{match, skip, limit, lookup}

	cursor, err := r.Collection.Aggregate(ctx, pipe)
	if err != nil {
		return nil, 0, err
	}

	err = cursor.All(ctx, &entity)
	if err != nil {
		return nil, 0, err
	}

	cPipe := []primitive.M{
		match,
		{
			"$group": primitive.M{
				"_id": nil,
				"counter": primitive.M{
					"$sum": 1,
				},
			},
		},
	}

	cursor2, err := r.Collection.Aggregate(ctx, cPipe)
	if err != nil {
		return nil, 0, err
	}

	err = cursor2.All(ctx, &counter)
	if err != nil {
		return nil, 0, err
	}

	if len(counter) <= 0 {
		return nil, 0, nil
	}

	return entity, counter[0].Counter, nil
}
