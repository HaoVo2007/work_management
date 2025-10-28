package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Columns struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	BoardID   string             `json:"board_id" bson:"board_id"`
	Name      string             `json:"name" bson:"name"`
	Color     string             `json:"color" bson:"color"`
	Position  int64              `json:"position" bson:"position"`
	CreatedBy string             `json:"created_by" bson:"created_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
