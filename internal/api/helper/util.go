package helper

import (
	"errors"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/storage"
)

func GetDomain(c *gin.Context) (string, error) {
	origin := c.Request.URL.Query().Get("origin")
	if origin != "" {
		return origin, nil
	}

	origin = c.Request.Header.Get("Origin")
	if origin == "" {
		return "", errors.New("origin header not found")
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
