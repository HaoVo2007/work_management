package boards

import (
	"context"
	"work-management/internal/domain/boards/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BoardRepository interface {
	CreateBoard(ctx context.Context, board *model.Boards) error
	GetAllBoards(ctx context.Context) ([]*model.Boards, error)
	GetBoardById(ctx context.Context, boardID primitive.ObjectID) (*model.Boards, error)
	UpdateBoard(ctx context.Context, id primitive.ObjectID, board *model.Boards) error
	DeleteBoard(ctx context.Context, boardID primitive.ObjectID) error
	GetBoardsByUserID(ctx context.Context, userID string) ([]*model.Boards, error)
}

type boardRepository struct {
	BoardRepository *mongo.Collection
}

func NewBoardRepository(collection *mongo.Collection) BoardRepository {
	return &boardRepository{
		BoardRepository: collection,
	}
}

func (r *boardRepository) CreateBoard(ctx context.Context, board *model.Boards) error {
	_, err := r.BoardRepository.InsertOne(ctx, board)
	return err
}

func (r *boardRepository) GetAllBoards(ctx context.Context) ([]*model.Boards, error) {

	var boards []*model.Boards

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

func (r *boardRepository) GetBoardById(ctx context.Context, boardID primitive.ObjectID) (*model.Boards, error) {

	var board model.Boards

	filter := bson.M{"_id": boardID}

	err := r.BoardRepository.FindOne(ctx, filter).Decode(&board)
	if err != nil {
		return nil, err
	}

	return &board, nil

}

func (r *boardRepository) UpdateBoard(ctx context.Context, id primitive.ObjectID, board *model.Boards) error {
	_, err := r.BoardRepository.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": board})
	return err
}

func (r *boardRepository) DeleteBoard(ctx context.Context, boardID primitive.ObjectID) error {
	_, err := r.BoardRepository.DeleteOne(ctx, bson.M{"_id": boardID})
	return err
}

func (r *boardRepository) GetBoardsByUserID(ctx context.Context, userID string) ([]*model.Boards, error) {

	var boards []*model.Boards

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
