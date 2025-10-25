package board

import (
	"context"
	"fmt"
	"work-management/internal/domain/board/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BoardRepository interface {
	Create(ctx context.Context, board *model.Board) error
	GetAll(ctx context.Context) ([]*model.Board, error)
	GetBoardById(ctx context.Context, boardID primitive.ObjectID) (*model.Board, error)
	GetBoardsByUserID(ctx context.Context, userID string) ([]*model.Board, error)
}

type boardRepository struct {
	BoardRepository *mongo.Collection
}

func NewBoardRepository(collection *mongo.Collection) BoardRepository {
	return &boardRepository{
		BoardRepository: collection,
	}
}

func (r *boardRepository) Create(ctx context.Context, board *model.Board) error {
	_, err := r.BoardRepository.InsertOne(ctx, board)
	return err
}

func (r *boardRepository) GetAll(ctx context.Context) ([]*model.Board, error) {

	var boards []*model.Board

	cursor, err := r.BoardRepository.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &boards); err != nil {
		return nil, err
	}

	return boards, nil
}

func (r *boardRepository) GetBoardById(ctx context.Context, boardID primitive.ObjectID) (*model.Board, error) {
	
	var board model.Board

	filter := bson.M{"_id": boardID}

	err := r.BoardRepository.FindOne(ctx, filter).Decode(&board)
	if err != nil {
		return nil, err
	}

	return &board, nil

}

func (r *boardRepository) GetBoardsByUserID(ctx context.Context, userID string) ([]*model.Board, error) {

	var boards []*model.Board

	filter := bson.M{"members": userID}

	cursor, err := r.BoardRepository.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &boards); err != nil {
		return nil, err
	}

	return boards, nil

}
