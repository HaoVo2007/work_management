package http

import (
	"work-management/configs"
	"work-management/internal/domain/boards"
	"work-management/internal/domain/columns"
	"work-management/internal/domain/tasks"
	"work-management/internal/domain/users"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(r *gin.Engine, db *mongo.Client, cfg *configs.Config) {

	userCollection := db.Database(cfg.MongoDB).Collection("users")
	boardCollection := db.Database(cfg.MongoDB).Collection("boards")
	columnCollection := db.Database(cfg.MongoDB).Collection("columns")
	taskCollection := db.Database(cfg.MongoDB).Collection("tasks")

	userRepo := users.NewRepository(userCollection)
	columnRepo := columns.NewColumnRepository(columnCollection)
	boardRepo := boards.NewBoardRepository(boardCollection)
	taskRepo := tasks.NewTaskRepository(taskCollection)

	userService := users.NewService(userRepo)
	boardService := boards.NewBoardService(boardRepo, columnRepo, userRepo, taskRepo)
	columnService := columns.NewColumnService(columnRepo, boardRepo)
	taskService := tasks.NewTaskService(taskRepo)

	users.NewHandler(r, userService)
	boards.NewBoardHandler(r, boardService)
	columns.NewColumnHandler(r, columnService)
	tasks.NewTaskHandler(r, taskService)
}
