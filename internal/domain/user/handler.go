package user

import "github.com/gin-gonic/gin"

type Handler struct {
	service Service
}

func NewHandler(r *gin.Engine, service Service) {
	
	handler := &Handler{
		service: service,
	}

	group := r.Group("api/v1") 
	{
		group.POST("/users", handler.CreateUser)
	}
	
}

func (h *Handler) CreateUser(c *gin.Context) {
	
}