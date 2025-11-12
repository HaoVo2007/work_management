package shared

import (
	"context"
	"work-management/internal/domain/boards/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardShared interface {
	GetBoardById(ctx context.Context, boardID primitive.ObjectID) (*model.Boards, error)
}
