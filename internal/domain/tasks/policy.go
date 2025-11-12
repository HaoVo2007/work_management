package tasks

import (
	"errors"
	"fmt"
	"work-management/internal/domain/columns/model"
)

type TaskPolicy struct{}

func NewTaskPolicy() *TaskPolicy {
	return &TaskPolicy{}
}

var ErrPermissionDenied = errors.New("permission denied")

func (p *TaskPolicy) CanCreateTask(column *model.Columns, userID string) error {
	if column.CreatedBy != userID {
		return fmt.Errorf("%w: only the creator of the column can create a task", ErrPermissionDenied)
	}
	return nil
}

func (p *TaskPolicy) CanUpdateTask(task *model.Columns, userID string) error {
	if task.CreatedBy != userID {
		return fmt.Errorf("%w: only the creator of the task can update it", ErrPermissionDenied)
	}
	return nil
}

func (p *TaskPolicy) CanDeleteTask(column *model.Columns, userID string) error {
	if column.CreatedBy != userID {
		return fmt.Errorf("%w: only the creator of the task can delete it", ErrPermissionDenied)
	}
	return nil
}
