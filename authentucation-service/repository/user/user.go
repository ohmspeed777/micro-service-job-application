package user

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
	collectionName = "users"
)

type Repo struct {
	repository.Repo
}

func NewRepo(db *mongo.Database) ports.IUserRepo {
	return &Repo{
		Repo: repository.Repo{
			Collection: db.Collection(collectionName),
		},
	}
}

func (r *Repo) FindOneByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	entity := &domain.User{}
	err := r.Collection.FindOne(ctx, primitive.M{"email": email}).Decode(entity)

	return entity, err
}
