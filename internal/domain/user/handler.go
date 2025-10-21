package user

import (
	"context"
	"fmt"
	"time"
	"work-management/internal/pkg/aws"
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

	group := r.Group("api/v1")
	{
		group.GET("/users", handler.GetUser)

		group.POST("/users/upload", handler.UploadUsers)
	}

}

func (h *Handler) GetUser(c *gin.Context) {

	key := c.Query("key")

	if key == "" {
		response.BadRequest(c, fmt.Errorf("missing key"))
		return
	}

	url, err := aws.GetPresignedURL(context.Background(), key, 15*time.Minute)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, "Avatar URL", url)
}

func (h *Handler) UploadUsers(c *gin.Context) {

	file, err := c.FormFile("avatar")
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	key, err := aws.UploadPrivateFile(context.Background(), file, "avatars")
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, "Avatar uploaded successfully", gin.H{"key": key})
}
