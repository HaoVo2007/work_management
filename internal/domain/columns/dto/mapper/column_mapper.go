package mapper

import (
	"work-management/internal/domain/columns/dto/response"
	"work-management/internal/domain/columns/model"
)

func ToColumnResponse(column *model.Columns) *response.ColumnResponse {
	return &response.ColumnResponse{
		ID:       column.ID.Hex(),
		Name:     column.Name,
		Position: int(column.Position),
	}
}

func ToColumnResponses(columns []*model.Columns) []*response.ColumnResponse {
	res := make([]*response.ColumnResponse, 0, len(columns))
	for _, c := range columns {
		res = append(res, ToColumnResponse(c))
	}
	return res
}
