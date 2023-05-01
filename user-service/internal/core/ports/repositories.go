package ports

import "app/internal/core/domain"

type IUserRepo interface {
	Create(i interface{}) error
	Update(i interface{}) error
	HardDelete(i interface{}) error
	FindOneByEmail(email string) (*domain.User, error)
	FindOneByID(id string) (*domain.User, error)
}


