package mapper

import (
	boardDTO "work-management/internal/domain/boards/dto/response"
	"work-management/internal/domain/boards/model"
	columnDTO "work-management/internal/domain/columns/dto/response"
	userDTO "work-management/internal/domain/users/dto/response"
)

func ToBoardResponse(
	board *model.Boards,
	user *userDTO.UserResponse,
	columns []*columnDTO.ColumnResponse,
	members []*userDTO.UserResponse,
) *boardDTO.BoardResponse {
	return &boardDTO.BoardResponse{
		ID:         board.ID.Hex(),
		Name:       board.Name,
		Background: *board.Background,
		Color:      *board.Color,
		Icon:       *board.Icon,
		Members:    append([]*userDTO.UserResponse{}, members...),
		Columns:    append([]*columnDTO.ColumnResponse{}, columns...),
		CreatedBy:  *user,
	}
}
