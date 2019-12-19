package api

import (
	"net/url"
	"strings"

	"github.com/systemli/ticker/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

//
func GetDomain(c *gin.Context) (string, error) {
	origin := c.Request.Header.Get("Origin")

	if origin == "" {
		origin = c.Request.URL.Query().Get("origin")
		if origin == "" {
			return "", errors.New("Origin header not found")
		} else {
			return origin, nil
		}
	}

	u, err := url.Parse(origin)
	if err != nil {
		return "", err
	}

	domain := strings.TrimPrefix(u.Host, "www.")
	if strings.Contains(domain, ":") {
		parts := strings.Split(domain, ":")
		domain = parts[0]
	}

	return domain, nil
}

func Me(c *gin.Context) (model.User, error) {
	var user model.User
	u, exists := c.Get(UserKey)
	if !exists {
		return user, errors.New(model.ErrorUserNotFound)
	}

	return u.(model.User), nil
}

func IsAdmin(c *gin.Context) bool {
	u, err := Me(c)
	if err != nil {
		return false
	}

	return u.IsSuperAdmin
}
