package comment

import "go.mongodb.org/mongo-driver/mongo"

type CommentRepository interface{}

type commentRepository struct {
	CommentCollection *mongo.Collection
}

func NewCommentRepository(collection *mongo.Collection) CommentRepository {
	return &commentRepository{
		CommentCollection: collection,
	}
}
