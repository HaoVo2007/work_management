package boards

import (
	"context"
	"fmt"
	"time"
	mapperBoard "work-management/internal/domain/boards/dto/mapper"
	"work-management/internal/domain/boards/dto/request"
	boardResponse "work-management/internal/domain/boards/dto/response"
	"work-management/internal/domain/boards/model"
	"work-management/internal/domain/columns"
	mapperColumn "work-management/internal/domain/columns/dto/mapper"
	"work-management/internal/domain/users"
	userDTO "work-management/internal/domain/users/dto/mapper"
	"work-management/internal/domain/users/dto/response"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardService interface {
	CreateBoard(ctx context.Context, req *request.CreateBoardRequest, userID string) (*model.Boards, error)
	GetAllBoards(ctx context.Context) ([]*boardResponse.BoardResponse, error)
	GetBoardById(ctx context.Context, boardID string) (*boardResponse.BoardResponse, error)
	UpdateBoard(ctx context.Context, boardID string, req *request.UpdateBoardRequest, userID string) (*model.Boards, error)
	DeleteBoard(ctx context.Context, boardID, userID string) error
	GetBoardsByUserID(ctx context.Context, userID string) ([]*model.Boards, error)
}

type boardService struct {
	BoardRepository  BoardRepository
	ColumnRepository columns.ColumnRepository
	UserRepository   users.Repository
}

func NewBoardService(
	boardRepository BoardRepository,
	columnRepository columns.ColumnRepository,
	userRepository users.Repository) BoardService {
	return &boardService{
		BoardRepository:  boardRepository,
		ColumnRepository: columnRepository,
		UserRepository:   userRepository,
	}
}

func (s *boardService) CreateBoard(ctx context.Context, req *request.CreateBoardRequest, userID string) (*model.Boards, error) {

	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	data := &model.Boards{
		ID:         primitive.NewObjectID(),
		Name:       req.Name,
		Background: req.Background,
		Color:      req.Color,
		Icon:       req.Icon,
		Members:    []string{userID},
		CreatedBy:  userID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := s.BoardRepository.CreateBoard(ctx, data)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (s *boardService) GetAllBoards(ctx context.Context) ([]*boardResponse.BoardResponse, error) {

	boards, err := s.BoardRepository.GetAllBoards(ctx)
	if err != nil {
		return nil, err
	}

	var result []*boardResponse.BoardResponse

	for _, board := range boards {

		columnsModel, err := s.ColumnRepository.GetColumnsByBoardID(ctx, board.ID.Hex())
		if err != nil {
			return nil, err
		}
		columnResponses := mapperColumn.ToColumnResponses(columnsModel)

		userIDobj, err := primitive.ObjectIDFromHex(board.CreatedBy)
		if err != nil {
			return nil, err
		}

		user, err := s.UserRepository.FindByID(ctx, userIDobj)
		if err != nil {
			return nil, err
		}

		if user == nil {
			return nil, fmt.Errorf("creator not found")
		}

		userResponse := userDTO.ToUserResponse(user)

		members := make([]*response.UserResponse, 0)
		for _, m := range board.Members {
			objectID, err := primitive.ObjectIDFromHex(m)
			if err != nil {
				return nil, err
			}
			memberUser, _ := s.UserRepository.FindByID(ctx, objectID)
			if memberUser != nil {
				members = append(members, userDTO.ToUserResponse(memberUser))
			}
		}

		result = append(result, mapperBoard.ToBoardResponse(board, userResponse, columnResponses, members))
	}

	return result, nil
}

func (s *boardService) GetBoardById(ctx context.Context, boardID string) (*boardResponse.BoardResponse, error) {

	if boardID == "" {
		return nil, fmt.Errorf("board_id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		return nil, err
	}

	board, err := s.BoardRepository.GetBoardById(ctx, objectID)
	if err != nil {
		return nil, err
	}
	if board == nil {
		return nil, fmt.Errorf("board not found")
	}

	columnsModel, err := s.ColumnRepository.GetColumnsByBoardID(ctx, boardID)
	if err != nil {
		return nil, err
	}
	columnResponses := mapperColumn.ToColumnResponses(columnsModel)

	userIDobj, err := primitive.ObjectIDFromHex(board.CreatedBy)
	if err != nil {
		return nil, err
	}

	user, err := s.UserRepository.FindByID(ctx, userIDobj)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("creator not found")
	}

	userResponse := userDTO.ToUserResponse(user)

	members := make([]*response.UserResponse, 0)
	for _, m := range board.Members {
		objectID, err := primitive.ObjectIDFromHex(m)
		if err != nil {
			return nil, err
		}
		memberUser, _ := s.UserRepository.FindByID(ctx, objectID)
		if memberUser != nil {
			members = append(members, userDTO.ToUserResponse(memberUser))
		}
	}

	boardResponse := mapperBoard.ToBoardResponse(board, userResponse, columnResponses, members)

	return boardResponse, nil
}

func (s *boardService) UpdateBoard(ctx context.Context, boardID string, req *request.UpdateBoardRequest, userID string) (*model.Boards, error) {

	if boardID == "" {
		return nil, fmt.Errorf("board_id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		return nil, err
	}

	board, err := s.BoardRepository.GetBoardById(ctx, objectID)
	if err != nil {
		return nil, err
	}

	policy := NewBoardPolicy()
	err = policy.CanUpdateBoard(board, userID)
	if err != nil {
		return nil, err
	}

	if board == nil {
		return nil, fmt.Errorf("board not found")
	}

	if req.Name != "" {
		board.Name = req.Name
	}

	if req.Color != nil {
		board.Color = req.Color
	}

	if req.Background != nil {
		board.Background = req.Background
	}

	if req.Icon != nil {
		board.Icon = req.Icon
	}

	board.UpdatedAt = time.Now()

	err = s.BoardRepository.UpdateBoard(ctx, objectID, board)
	if err != nil {
		return nil, err
	}

	return board, nil

}

func (s *boardService) DeleteBoard(ctx context.Context, boardID, userID string) error {

	if boardID == "" {
		return fmt.Errorf("board_id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		return err
	}

	board, err := s.BoardRepository.GetBoardById(ctx, objectID)
	if err != nil {
		return err
	}

	if board == nil {
		return fmt.Errorf("board not found")
	}

	policy := NewBoardPolicy()
	err = policy.CanDeleteBoard(board, userID)
	if err != nil {
		return err
	}

	columns, err := s.ColumnRepository.GetColumnsByBoardID(ctx, boardID)
	if err != nil {
		return err
	}

	for _, column := range columns {
		err = s.ColumnRepository.DeleteColumn(ctx, column.ID)
		if err != nil {
			return err
		}
	}

	return s.BoardRepository.DeleteBoard(ctx, objectID)

}

func (s *boardService) GetBoardsByUserID(ctx context.Context, userID string) ([]*model.Boards, error) {

	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	return s.BoardRepository.GetBoardsByUserID(ctx, userID)

}
