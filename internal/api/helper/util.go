package helper

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/storage"
)

func GetOrigin(c *gin.Context) (string, error) {
	var origin string
	if c.Request.URL.Query().Has("origin") {
		origin = c.Request.URL.Query().Get("origin")
	} else {
		origin = c.Request.Header.Get("Origin")
	}

	u, err := url.Parse(origin)
	if err != nil {
		return "", err
	}

	if u.Scheme == "" || u.Host == "" {
		return "", errors.New("invalid origin")
	}

	return fmt.Sprintf("%s://%s", u.Scheme, u.Host), nil
}

func GetOriginHost(c *gin.Context) string {
	origin, err := GetOrigin(c)
	if err != nil {
		return "unknown"
	}

	parts := strings.Split(origin, "://")
	if len(parts) < 2 {
		return "unknown"
	}

	return parts[1]
}

func Me(c *gin.Context) (storage.User, error) {
	var user storage.User
	u, exists := c.Get("me")
	if !exists {
		return user, errors.New("me not found")
	}

	return u.(storage.User), nil
}

func IsAdmin(c *gin.Context) bool {
	u, err := Me(c)
	if err != nil {
		return false
	}

	return u.IsSuperAdmin
}

func Ticker(c *gin.Context) (storage.Ticker, error) {
	ticker, exists := c.Get("ticker")
	if !exists {
		return storage.Ticker{}, errors.New("ticker not found")
	}

	return ticker.(storage.Ticker), nil
}

func Message(c *gin.Context) (storage.Message, error) {
	message, exists := c.Get("message")
	if !exists {
		return storage.Message{}, errors.New("message not found")
	}

	return message.(storage.Message), nil
}

func User(c *gin.Context) (storage.User, error) {
	user, exists := c.Get("user")
	if !exists {
		return storage.User{}, errors.New("user not found")
	}

	return user.(storage.User), nil
}
