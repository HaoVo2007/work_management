package columns

import (
	"errors"
	"fmt"
	"work-management/internal/domain/boards/model"
	columnModel "work-management/internal/domain/columns/model"
)

type ColumPolicy struct{}

func NewColumnPolicy() *ColumPolicy {
	return &ColumPolicy{}
}

var ErrPermissionDenied = errors.New("permission denied")

func (p *ColumPolicy) CanCreateColumn(board *model.Boards, userID string) error {
	if board.CreatedBy != userID {
		return fmt.Errorf("%w: only the creator of the board can create a column", ErrPermissionDenied)
	}
	return nil
}

func (p *ColumPolicy) CanUpdateColumn(column *columnModel.Columns, userID string) error {
	if column.CreatedBy != userID {
		return fmt.Errorf("%w: only the creator of the column can update it", ErrPermissionDenied)
	}
	return nil
}

func (p *ColumPolicy) CanDeleteColumn(column *columnModel.Columns, userID string) error {
	if column.CreatedBy != userID {
		return fmt.Errorf("%w: only the creator of the column can delete it", ErrPermissionDenied)
	}
	return nil
}
