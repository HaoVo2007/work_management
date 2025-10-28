package columns

import (
	"context"
	"work-management/internal/domain/columns/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ColumnRepository interface {
	CreateColumn(ctx context.Context, data *model.Columns) error
	GetColumnByID(ctx context.Context, id primitive.ObjectID) (*model.Columns, error)
	GetColumnsByBoardID(ctx context.Context, boardID string) ([]*model.Columns, error)
	UpdateColumn(ctx context.Context, id primitive.ObjectID, data *model.Columns) error
	UpdatePosition(ctx context.Context, id primitive.ObjectID, position int) error
	DeleteColumn(ctx context.Context, id primitive.ObjectID) error
	GetMaxPositionByBoardID(ctx context.Context, boardID string) (int, error)
}

type columnRepository struct {
	columnCollection *mongo.Collection
}

func NewColumnRepository(collection *mongo.Collection) ColumnRepository {
	return &columnRepository{
		columnCollection: collection,
	}
}

func (r *columnRepository) CreateColumn(ctx context.Context, data *model.Columns) error {
	_, err := r.columnCollection.InsertOne(ctx, data)
	return err
}

func (r *columnRepository) GetColumnByID(ctx context.Context, id primitive.ObjectID) (*model.Columns, error) {
	var column model.Columns

	err := r.columnCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&column)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &column, err
}

func (r *columnRepository) GetColumnsByBoardID(ctx context.Context, boardID string) ([]*model.Columns, error) {

	var columns []*model.Columns

	cursor, err := r.columnCollection.Find(ctx, bson.M{"board_id": boardID})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &columns); err != nil {
		return nil, err
	}

	return columns, nil
	
}

func (r *columnRepository) UpdateColumn(ctx context.Context, id primitive.ObjectID, data *model.Columns) error {
	_, err := r.columnCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})
	return err
}

func (r *columnRepository) UpdatePosition(ctx context.Context, id primitive.ObjectID, position int) error {
	_, err := r.columnCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"position": position}})
	return err
}

func (r *columnRepository) DeleteColumn(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.columnCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *columnRepository) GetMaxPositionByBoardID(ctx context.Context, boardID string) (int, error) {

	var result struct {
		Position int `bson:"position"`
	}

	err := r.columnCollection.
		FindOne(ctx, bson.M{"board_id": boardID}, options.FindOne().SetSort(bson.D{{Key: "position", Value: -1}})).
		Decode(&result)
		
	if err == mongo.ErrNoDocuments {
		return 0, nil
	}

	if err != nil {
		return 0, err	
	}

	return result.Position, nil
}
