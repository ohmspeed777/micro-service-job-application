package ports

import "app/internal/core/domain"

type IProductRepo interface {
	Create(i interface{}) error
	Update(i interface{}) error
	HardDelete(i interface{}) error
	FindAll(query domain.Query) ([]*domain.Product, int64, error)
	FindOneByID(id string) (*domain.Product, error)
}
