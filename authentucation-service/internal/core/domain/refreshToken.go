package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type TokenType int

const (
	ACCESS_TOKEN = iota + 1
	REFRESH_TOKEN
)

type RefreshToken struct {
	Model  `bson:",inline"`
	UUID   string             `json:"-" bson:"uuid"`
	UserID primitive.ObjectID `json:"-" bson:"user_id"`
	Type   TokenType          `json:"-" bson:"type"`
}
