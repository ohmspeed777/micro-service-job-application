package order

import (
	"app/internal/core/domain"
)

type GetAllReq struct {
	domain.Query
}

type CreateReq struct {
	Items  []*domain.OrderItem `json:"items"`
}
