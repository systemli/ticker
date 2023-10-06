package me

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
)

func MeMiddleware(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.UserIdentifierMissing))
			return
		}

		user, err := storage.FindUserByID(int(userID.(float64)))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.UserNotFound))
			return
		}

		c.Set("me", user)
	}
}
