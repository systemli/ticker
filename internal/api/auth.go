package api

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/storage"
)

//UserKey represents the key for gin context holding the user
const UserKey = "user"

const identityKey = "id"

//AuthMiddleware returns the Middleware for authenticating and authorising users with JWT
func AuthMiddleware() *jwt.GinJWTMiddleware {
	m := &jwt.GinJWTMiddleware{
		Realm:         "ticker admin",
		Key:           []byte(Config.Secret),
		Timeout:       time.Hour * 24,
		MaxRefresh:    time.Hour * 24,
		Authenticator: Authenticator,
		Authorizator:  Authorizator,
		Unauthorized:  Unauthorized,
		PayloadFunc:   FillClaim,
		TimeFunc:      time.Now,
		TokenLookup:   "header: Authorization",
		IdentityKey:   identityKey,
	}

	mw, err := jwt.New(m)
	if err != nil {
		panic(err)
	}

	return mw
}

//
func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(identityKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserIdentifierMissing))
			return
		}

		var user User
		err := DB.One("ID", int(userID.(float64)), &user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
			return
		}

		c.Set(UserKey, user)
	}
}

//Authenticator returns the user and the possible authentication error.
func Authenticator(c *gin.Context) (interface{}, error) {
	type login struct {
		Username string `form:"username" json:"username" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	var form login
	if err := c.ShouldBind(&form); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	return UserAuthenticate(form.Username, form.Password)
}

//Authorizator returns true when the user is authorized.
func Authorizator(data interface{}, c *gin.Context) bool {
	id := int(data.(float64))

	user, err := FindUserByID(id)
	if err != nil {
		return false
	}

	return user.ID != 0
}

//
func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, NewJSONErrorResponse(ErrorCodeCredentials, message))
}

//
func FillClaim(data interface{}) jwt.MapClaims {
	if u, ok := data.(*User); ok {
		return jwt.MapClaims{
			identityKey: u.ID,
			"email":     u.Email,
			"roles":     roles(u),
		}
	}

	return jwt.MapClaims{}
}

func roles(u *User) []string {
	roles := []string{"user"}

	if u.IsSuperAdmin {
		roles = append(roles, "admin")
	}

	return roles
}
