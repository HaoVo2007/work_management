package response

import "time"

type CreateUserResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name" bson:"name"`
	Avatar       *string   `json:"avatar" bson:"avatar"`
	Email        string    `json:"email" bson:"email"`
	Password     string    `json:"password" bson:"password"`
	Role         string   `json:"role" bson:"role"`
	Token        string    `json:"token" bson:"token"`
	RefreshToken string    `json:"refresh_token" bson:"refresh_token"`
	InviteLink   string    `json:"invite_link" bson:"invite_link"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}
