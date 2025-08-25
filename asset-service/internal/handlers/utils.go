package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ExtractUserID(c *gin.Context) (uuid.UUID, error) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return uuid.Nil, ErrUserNotAuthenticated
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID format"})
		return uuid.Nil, ErrInvalidUserIDFormat
	}

	return userID, nil
}

var (
	ErrUserNotAuthenticated = fmt.Errorf("user not authenticated")
	ErrInvalidUserIDFormat  = fmt.Errorf("invalid user ID format")
)
