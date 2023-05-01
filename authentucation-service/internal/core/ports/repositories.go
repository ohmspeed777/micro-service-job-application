package ports

import "app/internal/core/domain"

type IUserRepo interface {
	Create(i interface{}) error
	Update(i interface{}) error
	HardDelete(i interface{}) error
	FindOneByEmail(email string) (*domain.User, error)
}


type IRefreshTokenRepo interface {
	Create(i interface{}) error
	Update(i interface{}) error
	HardDelete(i interface{}) error
}