package shared

import (
	"context"
	"work-management/internal/domain/tasks/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskShared interface {
	DeleteTask(ctx context.Context, taskID primitive.ObjectID) error
	GetTasksByColumnID(ctx context.Context, columnID string) ([]*model.Tasks, error)
}
