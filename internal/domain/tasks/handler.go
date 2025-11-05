package tasks

import (
	"context"
	"fmt"
	"work-management/internal/app/http/middleware"
	"work-management/internal/domain/tasks/dto/request"
	"work-management/internal/pkg/constants"
	"work-management/internal/pkg/response"

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

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req request.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		response.Unauthorized(c, fmt.Errorf("missing token"))
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	userID, exists := c.Get(constants.UserID)
	if !exists {
		response.Unauthorized(c, fmt.Errorf("missing user_id in token"))
		return
	}

	task, err := h.TaskService.CreateTask(ctx, &req, userID.(string))
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, "Task created successfully", task)
}
