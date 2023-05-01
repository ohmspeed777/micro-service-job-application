package product

import (
	"app/internal/core/domain"
	"app/internal/core/ports"
	"app/repository"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "products"
)

type Repo struct {
	repository.Repo
}

func NewRepo(db *mongo.Database) ports.IProductRepo {
	return &Repo{
		Repo: repository.Repo{
			Collection: db.Collection(collectionName),
		},
	}
}

func (r *Repo) FindAll(query domain.Query) ([]*domain.Product, int64, error) {
	data := []*domain.Product{}

	filter := primitive.M{}


	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	opt := options.Find().SetSkip(query.GetSkip()).SetLimit(query.GetLimit())

	cursor, err := r.Collection.Find(ctx, filter, opt)
	if err != nil {
		return nil, 0, err
	}

	err = cursor.All(ctx, &data)
	if err != nil {
		return nil, 0, err
	}

	counter, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return data, counter, nil
}


func (r *Repo) FindOneByID(id string) (*domain.Product, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	entity := &domain.Product{}
	err = r.Collection.FindOne(ctx, primitive.M{"_id": _id}).Decode(entity)

	return entity, err
}
