package api

import (
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
)

const UserKey = "user"

//
func AuthMiddleware() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:         "ticker admin",
		Key:           []byte(Config.Secret),
		Timeout:       time.Hour * 24,
		MaxRefresh:    time.Hour * 24,
		Authenticator: Authenticator,
		Authorizator:  Authorizator,
		Unauthorized:  Unauthorized,
		TimeFunc:      time.Now,
	}
}

//
func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.GetString("userID")

		if email == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserIdentifierMissing))
			return
		}

		var user User
		err := DB.One("Email", email, &user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
			return
		}

		c.Set(UserKey, user)
	}
}

//
func Authenticator(userID string, password string, c *gin.Context) (interface{}, bool) {
	return UserAuthenticate(userID, password)
}

//
func Authorizator(data interface{}, c *gin.Context) bool {
	return UserExists(data.(string))
}

//
func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, NewJSONErrorResponse(ErrorCodeCredentials, message))
}
