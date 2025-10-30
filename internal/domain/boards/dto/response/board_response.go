package response

import (
	columnDTO "work-management/internal/domain/columns/dto/response"
	userDTO "work-management/internal/domain/users/dto/response"
)

type BoardResponse struct {
	ID         string                      `json:"id"`
	Name       string                      `json:"name"`
	Background string                      `json:"background" bson:"background"`
	Color      string                      `json:"color" bson:"color"`
	Icon       string                      `json:"icon" bson:"icon"`
	Members    []*userDTO.UserResponse     `json:"members" bson:"members"`
	Columns    []*columnDTO.ColumnResponse `json:"columns" bson:"columns"`
	CreatedBy  userDTO.UserResponse        `json:"created_by" bson:"created_by"`
}
