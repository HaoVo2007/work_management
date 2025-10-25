package board

import (
	"context"
	"fmt"
	"work-management/internal/app/http/middleware"
	"work-management/internal/domain/board/dto/request"
	"work-management/internal/pkg/constants"
	"work-management/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type BoardHandler struct {
	BoardService BoardService
}

func NewBoardHandler(r *gin.Engine, boardService BoardService) {

	handler := &BoardHandler{
		BoardService: boardService,
	}

	api := r.Group("api/v1")
	{
		public := api.Group("/boards")
		public.Use(middleware.JWTAuthMiddleware())
		{
			public.POST("", handler.CreateBoard)
			public.GET("", handler.GetAllBoards)
			public.GET("/:id", handler.GetBoardById)
			// public.PUT("/:id", handler.UpdateBoard)
			// public.DELETE("/:id", handler.DeleteBoard)
			public.GET("/user", handler.GetBoardsByUserID)
		}
	}
}

func (h *BoardHandler) CreateBoard(c *gin.Context) {

	userID, exists := c.Get(constants.UserID)
	if !exists {
		response.Unauthorized(c, fmt.Errorf("missing user_id in token"))
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		response.Unauthorized(c, fmt.Errorf("missing token"))
		return
	}

	var req request.CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	board, err := h.BoardService.CreateBoard(ctx, &req, userID.(string))
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, "Board created successfully", board)

}

func (h *BoardHandler) GetAllBoards(c *gin.Context) {

	token := c.GetHeader("Authorization")
	if token == "" {
		response.Unauthorized(c, fmt.Errorf("missing token"))
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	boards, err := h.BoardService.GetAllBoards(ctx)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, "Boards retrieved successfully", boards)

}

func (h *BoardHandler) GetBoardById(c *gin.Context) {

	token := c.GetHeader("Authorization")
	if token == "" {
		response.Unauthorized(c, fmt.Errorf("missing token"))
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	boardID := c.Param("id")
	board, err := h.BoardService.GetBoardById(ctx, boardID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, "Board retrieved successfully", board)

}

func (h *BoardHandler) GetBoardsByUserID(c *gin.Context) {

	userID, exists := c.Get(constants.UserID)
	if !exists {
		response.Unauthorized(c, fmt.Errorf("missing user_id in token"))
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		response.Unauthorized(c, fmt.Errorf("missing token"))
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	boards, err := h.BoardService.GetBoardsByUserID(ctx, userID.(string))
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, "Boards retrieved successfully", boards)
	
}