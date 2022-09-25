package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
)

func NeedAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := helper.Me(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.UserIdentifierMissing))
			return
		}

		if !user.IsSuperAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorResponse(response.CodeDefault, response.UserIdentifierMissing))
			return
		}
	}
}
