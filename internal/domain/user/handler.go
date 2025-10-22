package user

import (
	"fmt"
	"work-management/internal/app/http/middleware"
	"work-management/internal/domain/user/dto/request"
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
			public.POST("/register", handler.CreateUser)
			public.POST("/login", handler.LoginUser)
		}

		auth := api.Group("users")
		auth.Use(middleware.JWTAuthMiddleware())
		{
			auth.POST("/logout", handler.LogoutUser)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuthMiddleware(), middleware.IsAdminMiddleware())
		{
			
		}
	}

}

func (h *Handler) CreateUser(c *gin.Context) {

	var req request.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	user, err := h.service.CreateUser(c, req)
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

	userID, exists := c.Get("user_id")
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
