package api

import (
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

	"git.codecoop.org/systemli/ticker/internal/model"
)

func AuthMiddleware() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: Authenticator,
		Authorizator: Authorizator,
		Unauthorized: Unauthorized,
		TimeFunc: time.Now,
	}
}

//
func Authenticator(userID string, password string, c *gin.Context) (string, bool) {
	//TODO: Implement User Lookup
	if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
		return userID, true
	}

	return userID, false
}

//
func Authorizator(userID string, c *gin.Context) bool {
	//TODO: Implement User Lookup
	if userID == "admin" {
		return true
	}

	return false
}

//
func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, model.NewJSONErrorResponse(model.ErrorCredentials, message))
}
