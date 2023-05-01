package domain

import (
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
}

type Query struct {
	Limit uint `json:"limit" path:"limit" query:"limit"`
	Page  uint `json:"page" path:"page" query:"page"`
}

func (q *Query) GetLimit() int64 {
	if q.Limit == 0 {
		q.Limit = 20
	}

	return int64(q.Limit)
}

func (q *Query) GetPage() int64 {
	if q.Page == 0 {
		q.Page = 1
	}

	return int64(q.Page)
}

func (q *Query) GetSkip() int64 {
	if q.Limit == 0 {
		q.Limit = 20
	}

	if q.Page == 0 {
		q.Page = 1
	}

	return int64((q.Page - 1) * q.Limit)
}

type ModelInterface interface {
	GetID() primitive.ObjectID
	SetID(id primitive.ObjectID)
	Stamp()
	UpdateStamp()
}

// SetID set id
func (model *Model) SetID(id primitive.ObjectID) {
	model.ID = id
}

// GetID get id
func (model *Model) GetID() primitive.ObjectID {
	return model.ID
}

// Stamp current time to model
func (model *Model) Stamp() {
	timeNow := time.Now()
	model.UpdatedAt = timeNow
	model.CreatedAt = timeNow
}

// UpdateStamp current updated at model
func (model *Model) UpdateStamp() {
	timeNow := time.Now()
	model.UpdatedAt = timeNow
	if model.CreatedAt.IsZero() {
		model.CreatedAt = timeNow
	}
}

type CommonResponse struct {
	Message string `json:"message"`
}

func NewOkMessage() *CommonResponse {
	return &CommonResponse{
		Message: "ok",
	}
}

type GetOneReq struct {
	ID string `json:"id" path:"id" query:"id" param:"id"`
}

type PageInfo struct {
	Page          int64 `json:"page"`
	Size          int64 `json:"size"`
	NumOfPages    int64 `json:"num_of_pages"`
	AllOFEntities int64 `json:"all_of_entities"`
}

type ResponseInfo struct {
	PageInfo PageInfo `json:"page_info"`
	Entities any      `json:"entities"`
}

func NewPagination(entities any, page, size, counter int64) *ResponseInfo {
	res := &ResponseInfo{}
	res.Entities = entities
	res.PageInfo.NumOfPages = int64(math.Ceil(float64(counter) / float64(size)))
	res.PageInfo.Page = page
	res.PageInfo.Size = size
	res.PageInfo.AllOFEntities = counter
	return res
}

type Counter struct {
	Counter int64 `bson:"counter"`
}
