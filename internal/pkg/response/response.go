package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: "Bad Request",
		Error:   err.Error(),
	})
}

func Unauthorized(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Message: "Unauthorized",
		Error:   err.Error(),
	})
}

func Forbidden(c *gin.Context, err error) {
	c.JSON(http.StatusForbidden, APIResponse{
		Success: false,
		Message: "Forbidden",
		Error:   err.Error(),
	})
}

func NotFound(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Message: "Not Found",
		Error:   err.Error(),
	})
}

func InternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Message: "Internal Server Error",
		Error:   err.Error(),
	})
}
