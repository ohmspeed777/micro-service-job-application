package product

import "app/internal/core/domain"

type GetAllReq struct {
	domain.Query
}

type CreateReq struct {
	Price float64 `json:"price"`
	Stock uint    `json:"stock"`
	Name  string  `json:"name"`
}