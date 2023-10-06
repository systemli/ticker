package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
)

func PrefetchUser(s storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("userID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.UserIdentifierMissing))
			return
		}

		user, err := s.FindUserByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeNotFound, response.UserNotFound))
			return
		}

		c.Set("user", user)
	}
}
