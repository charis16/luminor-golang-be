package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Standard error response format
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
}

func RespondError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, ErrorResponse{
		Message: message,
		Code:    code,
	})
}

// Standard success response (optional helper)
func RespondSuccess(c *gin.Context, data gin.H) {
	c.JSON(http.StatusOK, data)
}
