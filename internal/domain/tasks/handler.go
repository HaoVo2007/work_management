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
			public.GET("", handler.GetAllTasks)
			public.GET("/:id", handler.GetTaskByID)
			public.PUT("/:id", handler.UpdateTask)
			public.DELETE("/:id", handler.DeleteTask)
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

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		response.Unauthorized(c, fmt.Errorf("missing token"))
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	tasks, err := h.TaskService.GetAllTasks(ctx)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, "Tasks retrieved successfully", tasks)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, fmt.Errorf("missing task_id"))
		return
	}
	
	token := c.GetHeader("Authorization")
	if token == "" {
		response.Unauthorized(c, fmt.Errorf("missing token"))
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	task, err := h.TaskService.GetTaskByID(ctx, taskID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, "Task retrieved successfully", task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var req request.UpdateTaskRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		response.Unauthorized(c, fmt.Errorf("missing token"))
		return
	}

	file, err := c.FormFile("cover_photo")
	if err != nil && err.Error() != "http: no such file" {
		response.BadRequest(c, err)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, fmt.Errorf("missing task_id"))
		return
	}

	// Pass the separately retrieved file to the service
	task, err := h.TaskService.UpdateTask(ctx, taskID, &req, file)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, "Task updated successfully", task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, fmt.Errorf("missing task_id"))
		return
	}
	
	token := c.GetHeader("Authorization")
	if token == "" {
		response.Unauthorized(c, fmt.Errorf("missing token"))
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.TaskService.DeleteTask(ctx, taskID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, "Task deleted successfully", nil)
}