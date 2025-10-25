package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Board struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Background *string            `json:"background" bson:"background"`
	Color      *string            `json:"color" bson:"color"`
	Members    []string           `json:"members" bson:"members"`
	CreatedBy  string             `json:"created_by" bson:"created_by"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}
