package auth

import (
	"errors"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
)

var log = logrus.WithField("package", "auth")

func AuthMiddleware(s storage.TickerStorage, secret string) *jwt.GinJWTMiddleware {
	config := &jwt.GinJWTMiddleware{
		Realm:         "ticker admin",
		Key:           []byte(secret),
		Timeout:       time.Hour * 24,
		MaxRefresh:    time.Hour * 24,
		Authenticator: Authenticator(s),
		Authorizator:  Authorizator(s),
		Unauthorized:  Unauthorized,
		PayloadFunc:   FillClaim,
		TimeFunc:      time.Now,
		TokenLookup:   "header: Authorization",
		IdentityKey:   "id",
	}

	middleware, err := jwt.New(config)
	if err != nil {
		log.WithError(err).Fatal()
	}

	return middleware
}

func Authenticator(s storage.TickerStorage) func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		type login struct {
			Username string `form:"username" json:"username" binding:"required"`
			Password string `form:"password" json:"password" binding:"required"`
		}

		var form login
		if err := c.ShouldBind(&form); err != nil {
			return "", jwt.ErrMissingLoginValues
		}

		user, err := s.FindUserByEmail(form.Username)
		if err != nil {
			log.WithError(err).Debug("user not found")
			return "", err
		}

		if user.Authenticate(form.Password) {
			return user, nil
		}

		return "", errors.New("authentication failed")
	}
}

func Authorizator(s storage.TickerStorage) func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		id := int(data.(float64))

		user, err := s.FindUserByID(id)
		if err != nil {
			log.WithError(err).WithField("data", data).Debug("user not found")
		}

		return user.ID != 0
	}
}

func Unauthorized(c *gin.Context, code int, message string) {
	log.WithFields(logrus.Fields{"code": code, "message": message, "url": c.Request.URL.String()}).Debug("unauthorized")
	c.JSON(code, response.ErrorResponse(response.CodeBadCredentials, response.Unauthorized))
}

func FillClaim(data interface{}) jwt.MapClaims {
	if u, ok := data.(storage.User); ok {
		return jwt.MapClaims{
			"id":    u.ID,
			"email": u.Email,
			"roles": roles(u),
		}
	}

	return jwt.MapClaims{}
}

func roles(u storage.User) []string {
	roles := []string{"user"}

	if u.IsSuperAdmin {
		roles = append(roles, "admin")
	}

	return roles
}
