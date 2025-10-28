package mapper

import (
	"work-management/internal/domain/users/dto/response"
	"work-management/internal/domain/users/model"
)

func ToUserResponse(user *model.Users) *response.UserResponse {

	if user == nil {
		return nil
	}

	var avatar string
	if user.Avatar != nil {
		avatar = *user.Avatar
	}

	return &response.UserResponse{
		ID:     user.ID.Hex(),
		Name:   user.Name,
		Email:  user.Email, 
		Avatar: avatar,
	}

}
