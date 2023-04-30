package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
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