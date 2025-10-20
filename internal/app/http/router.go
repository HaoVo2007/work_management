package http

import (
	"work-management/configs"
	"work-management/internal/domain/user"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(r *gin.Engine, db *mongo.Client, cfg *configs.Config) {

	userCollection := db.Database(cfg.MongoDB).Collection("users")
	userRepo := user.NewRepository(userCollection)
	userService := user.NewService(userRepo)
	user.NewHandler(r, userService)

}
