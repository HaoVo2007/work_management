package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comments struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	TaskID    string             `json:"task_id" bson:"task_id"`
	Content   *string            `json:"content" bson:"content"`
	File      *string            `json:"file" bson:"file"`
	Replies   []*Comments        `json:"replies" bson:"replies"`
	CreatedBy string             `json:"created_by" bson:"created_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
