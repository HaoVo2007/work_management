package tasks

import "go.mongodb.org/mongo-driver/mongo"

type TaskRepository interface {

}

type taskRepository struct {
	TaskCollection *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection) TaskRepository {
	return &taskRepository{
		TaskCollection: collection,
	}
}