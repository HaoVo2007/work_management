package users

import (
	"context"
	"work-management/internal/domain/users/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, user *model.Users) error
	UpdateByID(ctx context.Context, userID primitive.ObjectID, updateFields bson.M) error
	FindByEmail(ctx context.Context, email string) (*model.Users, error)
	FindByID(ctx context.Context, userID primitive.ObjectID) (*model.Users, error)	
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) Repository {
	return &repository{collection: collection}
}

func (r *repository) Create(ctx context.Context, user *model.Users) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *repository) UpdateByID(ctx context.Context, userID primitive.ObjectID, updateFields bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": updateFields})
	return err
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*model.Users, error) {

	var user model.Users

	filter := bson.M{"email": email}

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil

}

func (r *repository) FindByID(ctx context.Context, userID primitive.ObjectID) (*model.Users, error)	{

	var user model.Users

	filter := bson.M{"_id": userID}

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
	
}
