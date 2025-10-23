package user

import (
	"context"
	"fmt"
	"work-management/internal/app/http/middleware"
	"work-management/internal/domain/user/dto/request"
	"work-management/internal/pkg/constants"
	"work-management/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(r *gin.Engine, service Service) {

	handler := &Handler{
		service: service,
	}

	api := r.Group("api/v1")
	{
		public := api.Group("/users")
		{
			public.POST("/register", handler.RegisterUser)
			public.POST("/login", handler.LoginUser)
		}

		auth := api.Group("users")
		auth.Use(middleware.JWTAuthMiddleware())
		{
			auth.POST("/logout", handler.LogoutUser)
			auth.POST("/upload/avatar", handler.UploadAvatar)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuthMiddleware(), middleware.IsAdminMiddleware())
		{
			
		}
	}

}

func (h *Handler) RegisterUser(c *gin.Context) {

	var req request.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	user, err := h.service.RegisterUser(c, req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, "User created successfully", user)

}

func (h *Handler) LoginUser(c *gin.Context) {

	var req request.LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	token, err := h.service.LoginUser(c, req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, "User logged in successfully", token)

}

func (h *Handler) LogoutUser(c *gin.Context) {

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

	err := h.service.LogoutUser(c, userID.(string))
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, "User logged out successfully", nil)

}

func (h *Handler) UploadAvatar(c *gin.Context) {

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

	file, err := c.FormFile("avatar")
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	avatar, err := h.service.UploadAvatar(ctx, userID.(string), file)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, "Avatar uploaded successfully", avatar)

}
