package columns

import (
	"errors"
	"work-management/internal/domain/boards/model"
	columnModel "work-management/internal/domain/columns/model"
)

type ColumPolicy struct{}

func NewColumnPolicy() *ColumPolicy {
	return &ColumPolicy{}
}

func (p *ColumPolicy) CanCreateColumn(board *model.Boards, userID string) error {
	if board.CreatedBy != userID {
		return errors.New("permission denied: only the creator of the board can create a column")
	}
	return nil
}

func (p *ColumPolicy) CanUpdateColumn(column *columnModel.Columns, userID string) error {
	if column.CreatedBy != userID {
		return errors.New("permission denied: only the creator of the column can update it")
	}
	return nil
}

func (p *ColumPolicy) CanDeleteColumn(column *columnModel.Columns, userID string) error {
	if column.CreatedBy != userID {
		return errors.New("permission denied: only the creator of the column can update it")
	}
	return nil
}
