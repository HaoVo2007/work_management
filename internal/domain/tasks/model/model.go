package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tasks struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	BoardID     string             `json:"board_id" bson:"board_id"`
	ColumnID    string             `json:"column_id" bson:"column_id"`
	Title       string             `json:"title" bson:"title"`
	Description *string            `json:"description" bson:"description"`
	Assgine     string             `json:"assignee" bson:"assignee"`
	StartDate   time.Time          `json:"start_date" bson:"start_date"`
	EndDate     time.Time          `json:"end_date" bson:"end_date"`
	Priority    int64              `json:"priority" bson:"priority"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
