package mapper

import (
	"context"
	"work-management/internal/domain/columns/dto/response"
	"work-management/internal/domain/columns/model"
	"work-management/internal/domain/tasks"
)

func ToColumnResponse(column *model.Columns) *response.ColumnResponse {
	return &response.ColumnResponse{
		ID:       column.ID.Hex(),
		Name:     column.Name,
		Position: int(column.Position),
	}
}

func ToColumnResponses(ctx context.Context, columns []*model.Columns, taskRepository tasks.TaskRepository) ([]*response.ColumnResponse, error) {
	res := make([]*response.ColumnResponse, 0, len(columns))
	for _, c := range columns {
		tasksModel, err := taskRepository.GetTasksByColumnID(ctx, c.ID.Hex())
		if err != nil {
			return nil, err
		}

		taskResponses := make([]*response.TaskResponse, 0, len(tasksModel))
		for _, task := range tasksModel {
			taskResponses = append(taskResponses, &response.TaskResponse{
				ID:    task.ID.Hex(),
				Title: task.Name,
			})
		}
		res = append(res, &response.ColumnResponse{
			ID:       c.ID.Hex(),
			Name:     c.Name,
			Position: int(c.Position),
			Tasks:    taskResponses,
		})
	}
	return res, nil
}
