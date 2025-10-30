package boards

import (
	"errors"
	"fmt"
	"work-management/internal/domain/boards/model"
)

type BoardPolicy struct{}

func NewBoardPolicy() *BoardPolicy {
	return &BoardPolicy{}
}

var ErrPermissionDenied = errors.New("permission denied")

func (p *BoardPolicy) CanUpdateBoard(board *model.Boards, userID string) error {
	if board.CreatedBy != userID {
		return fmt.Errorf("%w: only the creator can update this board", ErrPermissionDenied)
	}
	return nil
}

func (p *BoardPolicy) CanDeleteBoard(board *model.Boards, userID string) error {
	if board.CreatedBy != userID {
		return fmt.Errorf("%w: only the creator can delete this board", ErrPermissionDenied)
	}
	return nil
}
