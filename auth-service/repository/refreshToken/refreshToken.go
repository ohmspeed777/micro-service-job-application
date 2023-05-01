package refreshtoken

import (
	"app/internal/core/ports"
	"app/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName = "refresh_tokens"
)

type Repo struct {
	repository.Repo
}

func NewRepo(db *mongo.Database) ports.IRefreshTokenRepo {
	return &Repo{
		Repo: repository.Repo{
			Collection: db.Collection(collectionName),
		},
	}
}
