package mapper

import (
	"context"
	"time"
	"work-management/internal/domain/users/dto/response"
	"work-management/internal/domain/users/model"
	"work-management/internal/pkg/aws"
)

func ToUserResponse(ctx context.Context, user *model.Users) *response.UserResponse {

	if user == nil {
		return nil
	}

	var avatar string
	if user.Avatar != nil {
		avatarUrl, err := aws.GetPresignedURL(ctx, *user.Avatar, 24*time.Hour)
		if err != nil {
			return nil
		}
		avatar = *avatarUrl
	}

	return &response.UserResponse{
		ID:     user.ID.Hex(),
		Name:   user.Name,
		Email:  user.Email,
		Avatar: avatar,
	}

}
