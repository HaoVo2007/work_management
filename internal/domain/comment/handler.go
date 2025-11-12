package comment

import (
	"work-management/internal/app/http/middleware"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	CommentService CommentService
}

func NewCommentHandler(r *gin.Engine, commentService CommentService) {
	handler := &CommentHandler{
		CommentService: commentService,
	}

	api := r.Group("api/v1")
	{
		public := api.Group("/comments")
		public.Use(middleware.JWTAuthMiddleware())
		{
			public.POST("", handler.CreateComment)
		}
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	
}

