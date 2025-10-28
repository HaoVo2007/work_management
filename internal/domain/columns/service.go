package columns

import (
	"context"
	"fmt"
	"sort"
	"time"
	"work-management/internal/domain/columns/dto/request"
	"work-management/internal/domain/columns/model"
	"work-management/internal/domain/shared"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ColumnService interface {
	CreateColumn(ctx context.Context, req *request.CreateColumnRequest, userID string) (*model.Columns, error)
	UpdateColumn(ctx context.Context, columnID string, req *request.UpdateColumnRequest, userID string) (*model.Columns, error)
	DeleteColumn(ctx context.Context, columnID string, userID string) error
}

type columnService struct {
	columnRepository ColumnRepository
	boardRepository  shared.BoardGetter
}

func NewColumnService(columnRepository ColumnRepository,
	boardRepository shared.BoardGetter) ColumnService {
	return &columnService{
		columnRepository: columnRepository,
		boardRepository:  boardRepository,
	}
}

func (s *columnService) CreateColumn(ctx context.Context, req *request.CreateColumnRequest, userID string) (*model.Columns, error) {

	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if req.BoardID == "" {
		return nil, fmt.Errorf("board_id is required")
	}

	if req.Color == "" {
		return nil, fmt.Errorf("color is required")
	}

	objectID, err := primitive.ObjectIDFromHex(req.BoardID)
	if err != nil {
		return nil, err
	}

	board, err := s.boardRepository.GetBoardById(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if board == nil {
		return nil, fmt.Errorf("board not found")
	}

	policy := NewColumnPolicy()
	err = policy.CanCreateColumn(board, userID)
	if err != nil {
		return nil, err
	}

	maxPos, err := s.columnRepository.GetMaxPositionByBoardID(ctx, req.BoardID)
	if err != nil {
		return nil, err
	}

	data := &model.Columns{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		BoardID:   req.BoardID,
		Color:     req.Color,
		Position:  int64(maxPos + 1),
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.columnRepository.CreateColumn(ctx, data)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (s *columnService) UpdateColumn(ctx context.Context, columnID string, req *request.UpdateColumnRequest, userID string) (*model.Columns, error) {

	if columnID == "" {
		return nil, fmt.Errorf("column_id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(columnID)
	if err != nil {
		return nil, err
	}

	column, err := s.columnRepository.GetColumnByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if column == nil {
		return nil, fmt.Errorf("column not found")
	}

	policy := NewColumnPolicy()
	err = policy.CanUpdateColumn(column, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		column.Name = req.Name
	}

	if req.Color != "" {
		column.Color = req.Color
	}

	column.UpdatedAt = time.Now()

	err = s.columnRepository.UpdateColumn(ctx, objectID, column)
	if err != nil {
		return nil, err
	}

	return column, nil

}

func (s *columnService) DeleteColumn(ctx context.Context, columnID string, userID string) error {

	if columnID == "" {
		return fmt.Errorf("column_id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(columnID)
	if err != nil {
		return err
	}

	column, err := s.columnRepository.GetColumnByID(ctx, objectID)
	if err != nil {
		return err
	}

	if column == nil {
		return fmt.Errorf("column not found")
	}

	policy := NewColumnPolicy()
	err = policy.CanDeleteColumn(column, userID)
	if err != nil {
		return err
	}

	err = s.columnRepository.DeleteColumn(ctx, objectID)
	if err != nil {
		return err
	}

	columns, err := s.columnRepository.GetColumnsByBoardID(ctx, column.BoardID)
	if err != nil {
		return err
	}

	sort.Slice(columns, func(i, j int) bool {
		return columns[i].Position < columns[j].Position
	})

	for idx, c := range columns {
		if c.Position != int64(idx) { 
			c.Position = int64(idx)
			if err := s.columnRepository.UpdatePosition(ctx, c.ID, idx); err != nil {
				return fmt.Errorf("failed to update position for column %s: %w", c.ID.Hex(), err)
			}
		}
	}

	return nil

}
