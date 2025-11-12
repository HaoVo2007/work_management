package tasks

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"
	"work-management/internal/domain/columns"
	"work-management/internal/domain/tasks/dto/mapper"
	"work-management/internal/domain/tasks/dto/request"
	"work-management/internal/domain/tasks/dto/response"
	"work-management/internal/domain/tasks/model"

	"log"
	"work-management/internal/pkg/aws"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService interface {
	CreateTask(ctx context.Context, req *request.CreateTaskRequest, userID string) (*model.Tasks, error)
	GetAllTasks(ctx context.Context) ([]*model.Tasks, error)
	GetTaskByID(ctx context.Context, taskID string) (*response.TaskResponse, error)
	UpdateTask(ctx context.Context, taskID string, req *request.UpdateTaskRequest, file *multipart.FileHeader) (*model.Tasks, error)
	DeleteTask(ctx context.Context, taskID string) error
}

type taskService struct {
	TaskRepository   TaskRepository
	ColumnRepository columns.ColumnRepository
}

func NewTaskService(
	taskRepository TaskRepository,
	columnRepository columns.ColumnRepository,
) TaskService {
	return &taskService{
		TaskRepository:   taskRepository,
		ColumnRepository: columnRepository,
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

	objectID, err := primitive.ObjectIDFromHex(req.ColumnID)
	if err != nil {
		return nil, err
	}

	column, err := s.ColumnRepository.GetColumnByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if column == nil {
		return nil, fmt.Errorf("column not found")
	}

	policy := NewTaskPolicy()
	err = policy.CanCreateTask(column, userID)
	if err != nil {
		return nil, err
	}

	data := &model.Tasks{
		ID:        primitive.NewObjectID(),
		Name:      req.Title,
		ColumnID:  req.ColumnID,
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.TaskRepository.CreateTask(ctx, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *taskService) GetAllTasks(ctx context.Context) ([]*model.Tasks, error) {
	return s.TaskRepository.GetAllTasks(ctx)
}

func (s *taskService) GetTaskByID(ctx context.Context, taskID string) (*response.TaskResponse, error) {
	if taskID == "" {
		return nil, fmt.Errorf("task_id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, err
	}

	task, err := s.TaskRepository.GetTaskByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return mapper.ToTaskResponse(ctx, task), nil
}

func (s *taskService) UpdateTask(ctx context.Context, taskID string, req *request.UpdateTaskRequest, file *multipart.FileHeader) (*model.Tasks, error) {

	if taskID == "" {
		return nil, fmt.Errorf("task_id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, err
	}

	task, err := s.TaskRepository.GetTaskByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, fmt.Errorf("task not found")
	}

	column, err := s.ColumnRepository.GetColumnByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if column == nil {
		return nil, fmt.Errorf("column not found")
	}

	policy := NewTaskPolicy()
	err = policy.CanUpdateTask(column, task.CreatedBy)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		task.Name = *req.Title
	}

	if file != nil {
		if task.CoverPhoto != nil && *task.CoverPhoto != "" {
			err = aws.DeleteFile(ctx, *task.CoverPhoto)
			if err != nil {
				log.Printf("Error deleting old task cover photo %s: %v", *task.CoverPhoto, err)
			}
		}

		key, err := aws.UploadPrivateFile(ctx, file, "tasks")
		if err != nil {
			return nil, fmt.Errorf("failed to upload new cover photo: %w", err)
		}
		task.CoverPhoto = &key
	}

	if req.Description != nil {
		task.Description = req.Description
	}

	if req.Assignee != nil {
		task.Assgine = *req.Assignee
	}

	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %w", err)
		}
		task.StartDate = startDate
	}

	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}
		task.EndDate = endDate
	}

	if req.Priority != nil {
		task.Priority = *req.Priority
	}

	task.UpdatedAt = time.Now()

	err = s.TaskRepository.UpdateTask(ctx, objectID, task)
	if err != nil {
		return nil, fmt.Errorf("failed to update task in repository: %w", err)
	}

	return task, nil
}

func (s *taskService) DeleteTask(ctx context.Context, taskID string) error {
	if taskID == "" {
		return fmt.Errorf("task_id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}

	task, err := s.TaskRepository.GetTaskByID(ctx, objectID)
	if err != nil {
		return err
	}

	if task == nil {
		return fmt.Errorf("task not found")
	}

	column, err := s.ColumnRepository.GetColumnByID(ctx, objectID)
	if err != nil {
		return err
	}

	if column == nil {
		return fmt.Errorf("column not found")
	}

	policy := NewTaskPolicy()
	err = policy.CanDeleteTask(column, task.CreatedBy)
	if err != nil {
		return err
	}

	err = s.TaskRepository.DeleteTask(ctx, objectID)
	if err != nil {
		return err
	}

	return nil
}
