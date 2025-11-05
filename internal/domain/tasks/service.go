package tasks

import (
	"context"
	"fmt"
	"time"
	"work-management/internal/domain/tasks/dto/request"
	"work-management/internal/domain/tasks/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService interface {
	CreateTask(ctx context.Context, req *request.CreateTaskRequest, userID string) (*model.Tasks, error)
}

type taskService struct {
	TaskRepository TaskRepository
}

func NewTaskService(taskRepository TaskRepository) TaskService {
	return &taskService{
		TaskRepository: taskRepository,
	}
}

func (s *taskService) CreateTask(ctx context.Context, req *request.CreateTaskRequest, userID string) (*model.Tasks, error) {
	
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	if req.ColumnID == "" {
		return nil, fmt.Errorf("column_id is required")
	}

	data := &model.Tasks{
		ID:        primitive.NewObjectID(),
		Name:      req.Title,
		ColumnID:  req.ColumnID,
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.TaskRepository.CreateTask(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
