package user

import "go.mongodb.org/mongo-driver/mongo"

type Repository interface {
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) Repository {
	return &repository{collection: collection}
}
