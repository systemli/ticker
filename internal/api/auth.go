package api

import (
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
)

func AuthMiddleware() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:         "test zone",
		Key:           []byte("secret key"),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		Authenticator: Authenticator,
		Authorizator:  Authorizator,
		Unauthorized:  Unauthorized,
		TimeFunc:      time.Now,
	}
}

//
func Authenticator(userID string, password string, c *gin.Context) (string, bool) {
	return UserAuthenticate(userID, password)
}

//
func Authorizator(userID string, c *gin.Context) bool {
	return UserExists(userID)
}

//
func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, NewJSONErrorResponse(ErrorCredentials, message))
}
