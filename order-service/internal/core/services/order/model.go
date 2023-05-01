package order

import (
	"app/internal/core/domain"
)

type GetAllReq struct {
	domain.Query
	User string
}

type CreateReq struct {
	Items []*domain.OrderItem `json:"items"`

	// mock up
	Total float64 `json:"total"`
}
