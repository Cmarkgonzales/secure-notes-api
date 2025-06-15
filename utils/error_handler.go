package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespondError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"error": message,
	})
}

func InternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "Internal Server Error",
		"details": err.Error(),
	})
}
