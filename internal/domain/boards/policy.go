package boards

import (
	"errors"
	"work-management/internal/domain/boards/model"
)

type BoardPolicy struct{}

func NewBoardPolicy() *BoardPolicy {
	return &BoardPolicy{}
}

func (p *BoardPolicy) CanUpdateBoard(board *model.Boards, userID string) error {
	if board.CreatedBy != userID {
		return errors.New("permission denied: only the creator can update this board")
	}
	return nil
}

func (p *BoardPolicy) CanDeleteBoard(board *model.Boards, userID string) error {
	if board.CreatedBy != userID {
		return errors.New("permission denied: only the creator can delete this board")
	}
	return nil
}
