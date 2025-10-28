package columns

import (
	"context"
	"fmt"
	"work-management/internal/app/http/middleware"
	"work-management/internal/domain/columns/dto/request"
	"work-management/internal/pkg/constants"
	"work-management/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type ColumnHandler struct {
	ColumnService ColumnService
}

func NewColumnHandler(r *gin.Engine, columnService ColumnService) {
	handler := &ColumnHandler{
		ColumnService: columnService,
	}

	api := r.Group("api/v1")
	{
		public := api.Group("/columns")
		public.Use(middleware.JWTAuthMiddleware())
		{
			public.POST("", handler.CreateColumn)
			// public.GET("", handler.GetAllColumns)
			// public.GET("/:id", handler.GetColumnByID)
			public.PUT("/:id", handler.UpdateColumn)
			public.DELETE("/:id", handler.DeleteColumn)
		}
	}
}

func (h *ColumnHandler) CreateColumn(c *gin.Context) {

	var req request.CreateColumnRequest

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

	column, err := h.ColumnService.CreateColumn(ctx, &req, userID.(string))
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, "Column created successfully", column)

}

func (h *ColumnHandler) UpdateColumn(c *gin.Context) {

	token := c.GetHeader("Authorization")
	if token == "" {
		response.Unauthorized(c, fmt.Errorf("missing token"))
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	columnID := c.Param("id")
	if columnID == "" {
		response.BadRequest(c, fmt.Errorf("missing column_id"))
		return
	}

	var req request.UpdateColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		response.Unauthorized(c, fmt.Errorf("missing user_id in token"))
		return
	}

	column, err := h.ColumnService.UpdateColumn(ctx, columnID, &req, userID.(string))
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, "Column updated successfully", column)

}

func (h *ColumnHandler) DeleteColumn(c *gin.Context) {

	token := c.GetHeader("Authorization")
	if token == "" {
		response.Unauthorized(c, fmt.Errorf("missing token"))
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	columnID := c.Param("id")
	if columnID == "" {
		response.BadRequest(c, fmt.Errorf("missing column_id"))
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		response.Unauthorized(c, fmt.Errorf("missing user_id in token"))
		return
	}

	err := h.ColumnService.DeleteColumn(ctx, columnID, userID.(string))
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, "Column deleted successfully", nil)

}
