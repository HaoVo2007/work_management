package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Avatar     *string             `json:"avatar" bson:"avatar"`
	Email      string             `json:"email" bson:"email"`
	Password   string             `json:"password" bson:"password"`
	Role       string             `json:"role" bson:"role"`
	InviteLink string             `json:"invite_link" bson:"invite_link"`
	CreatedAt  int64              `json:"created_at" bson:"created_at"`
	UpdatedAt  int64              `json:"updated_at" bson:"updated_at"`
}
