package ports

import "app/internal/core/domain"

type IProductRepo interface {
	Create(i interface{}) error
	Update(i interface{}) error
	HardDelete(i interface{}) error
	FindAll(query domain.Query) ([]*domain.Product, int64, error)
	FindOneByID(id string) (*domain.Product, error)
}

type IOrderRepo interface{
	Create(i interface{}) error
	Update(i interface{}) error
	FindOneByID(id string) (*domain.Order, error) 
	AggregateOneByID(id string) (*domain.OrderLookedUp, error)
	AggregateAllByUser(userId string, q domain.Query) ([]*domain.OrderLookedUp, int64, error)
}
