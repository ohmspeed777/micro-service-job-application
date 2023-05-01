package user

import (
	"app/internal/core/domain"
	"app/internal/core/ports"
	"app/repository/user"
	"context"
	"net/http"

	"github.com/ohmspeed777/go-pkg/corex"
	"github.com/ohmspeed777/go-pkg/errorx"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserService interface {
	GetMe(ctx context.Context) (*domain.User, error)
}

type Service struct {
	user ports.IUserRepo
}

func NewService(db *mongo.Database) IUserService {
	return &Service{
		user: user.NewRepo(db),
	}
}

func (s *Service) GetMe(ctx context.Context) (*domain.User, error) {
	c, err := corex.NewFromOutingContext(ctx)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not bind my context", err)
	}

	user, err := s.user.FindOneByID(c.User.ID)
	if err != nil {
		return nil, errorx.New(http.StatusBadRequest, "can not find my profile", err)
	}

	return user, nil
}
