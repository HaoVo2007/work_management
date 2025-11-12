package tasks

import (
	"context"
	"work-management/internal/domain/tasks/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, data *model.Tasks) error
	GetAllTasks(ctx context.Context) ([]*model.Tasks, error)
	GetTaskByID(ctx context.Context, taskID primitive.ObjectID) (*model.Tasks, error)
	GetTasksByColumnID(ctx context.Context, columnID string) ([]*model.Tasks, error)
	UpdateTask(ctx context.Context, taskID primitive.ObjectID, data *model.Tasks) error
	DeleteTask(ctx context.Context, taskID primitive.ObjectID) error
}

type taskRepository struct {
	TaskCollection *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection) TaskRepository {
	return &taskRepository{
		TaskCollection: collection,
	}
}

func (r *taskRepository) CreateTask(ctx context.Context, data *model.Tasks) error {
	_, err := r.TaskCollection.InsertOne(ctx, data)
	return err
}

func (r *taskRepository) GetAllTasks(ctx context.Context) ([]*model.Tasks, error) {
	filter := bson.M{}

	cursor, err := r.TaskCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*model.Tasks
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) GetTaskByID(ctx context.Context, taskID primitive.ObjectID) (*model.Tasks, error) {
	filter := bson.M{"_id": taskID}

	var task model.Tasks
	err := r.TaskCollection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &task, err
}

func (r *taskRepository) GetTasksByColumnID(ctx context.Context, columnID string) ([]*model.Tasks, error) {
	filter := bson.M{"column_id": columnID}

	cursor, err := r.TaskCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*model.Tasks
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, taskID primitive.ObjectID, data *model.Tasks) error {
	_, err := r.TaskCollection.UpdateOne(ctx, bson.M{"_id": taskID}, bson.M{"$set": data})
	return err
}

func (r *taskRepository) DeleteTask(ctx context.Context, taskID primitive.ObjectID) error {
	_, err := r.TaskCollection.DeleteOne(ctx, bson.M{"_id": taskID})
	return err
}
