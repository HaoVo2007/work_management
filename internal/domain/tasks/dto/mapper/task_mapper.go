package mapper

import (
	"context"
	"time"
	"work-management/internal/domain/tasks/dto/response"
	"work-management/internal/domain/tasks/model"
	"work-management/internal/pkg/aws"
)

func ToTaskResponse(ctx context.Context, task *model.Tasks) *response.TaskResponse {
	var coverPhoto *string
	if task.CoverPhoto != nil {
		coverPhotoUrl, err := aws.GetPresignedURL(ctx, *task.CoverPhoto, 24*time.Hour)
		if err != nil {
			return nil
		}
		coverPhoto = coverPhotoUrl
	}
	return &response.TaskResponse{
		ID:          task.ID.Hex(),
		ColumnID:    task.ColumnID,
		Name:        task.Name,
		CoverPhoto:  coverPhoto,
		Description: task.Description,
		Assgine:     task.Assgine,
		StartDate:   task.StartDate,
		EndDate:     task.EndDate,
		Priority:    task.Priority,
		CreatedBy:   task.CreatedBy,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
