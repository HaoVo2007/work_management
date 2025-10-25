package board

import (
	"context"
	"fmt"
	"time"
	"work-management/internal/domain/board/dto/request"
	"work-management/internal/domain/board/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardService interface {
	CreateBoard(ctx context.Context, req *request.CreateBoardRequest, userID string) (*model.Board, error)
	GetAllBoards(ctx context.Context) ([]*model.Board, error)
	GetBoardById(ctx context.Context, boardID string) (*model.Board, error)
	GetBoardsByUserID(ctx context.Context, userID string) ([]*model.Board, error)
}

type boardService struct {
	BoardRepository BoardRepository
}

func NewBoardService(boardRepository BoardRepository) BoardService {
	return &boardService{
		BoardRepository: boardRepository,
	}
}

func (s *boardService) CreateBoard(ctx context.Context, req *request.CreateBoardRequest, userID string) (*model.Board, error) {

	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	data := &model.Board{
		ID:         primitive.NewObjectID(),
		Name:       req.Name,
		Background: req.Background,
		Color:      req.Color,
		Members:    []string{userID},
		CreatedBy:  userID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := s.BoardRepository.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (s *boardService) GetAllBoards(ctx context.Context) ([]*model.Board, error) {
	return s.BoardRepository.GetAll(ctx)
}

func (s *boardService) GetBoardById(ctx context.Context, boardID string) (*model.Board, error) {

	if boardID == "" {
		return nil, fmt.Errorf("board_id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		return nil, err
	}

	return s.BoardRepository.GetBoardById(ctx, objectID)

}

func (s *boardService) GetBoardsByUserID(ctx context.Context, userID string) ([]*model.Board, error) {

	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	return s.BoardRepository.GetBoardsByUserID(ctx, userID)
	
}