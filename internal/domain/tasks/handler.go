package tasks

import (
	"work-management/internal/app/http/middleware"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	TaskService TaskService
}

func NewTaskHandler(r *gin.Engine, taskService TaskService) {
	handler := &TaskHandler{
		TaskService: taskService,
	}

	api := r.Group("api/v1")
	{
		public := api.Group("/tasks")
		public.Use(middleware.JWTAuthMiddleware())
		{
			public.POST("", handler.CreateTask)
			// public.GET("", handler.GetAllTasks)
			// public.GET("/:id", handler.GetTaskByID)
			// public.PUT("/:id", handler.UpdateTask)
			// public.DELETE("/:id", handler.DeleteTask)
		}
	}
}

func ( h *TaskHandler) CreateTask(c *gin.Context) {
	
}
