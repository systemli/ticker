package helper

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/storage"
)

type UtilTestSuite struct {
	suite.Suite
}

func (s *UtilTestSuite) TestGetOrigin() {
	s.Run("when origin is empty", func() {
		c := s.buildContext(url.URL{}, http.Header{})
		origin, err := GetOrigin(c)
		s.Equal("", origin)
		s.Equal("invalid origin", err.Error())
	})

	s.Run("when origin is not a valid URL", func() {
		c := s.buildContext(url.URL{}, http.Header{
			"Origin": []string{"localhost"},
		})
		origin, err := GetOrigin(c)
		s.Equal("", origin)
		s.Error(err)
	})

	s.Run("when origin is localhost", func() {
		c := s.buildContext(url.URL{}, http.Header{
			"Origin": []string{"http://localhost"},
		})
		origin, err := GetOrigin(c)
		s.Equal("http://localhost", origin)
		s.NoError(err)
	})

	s.Run("when origin is localhost with port", func() {
		c := s.buildContext(url.URL{}, http.Header{
			"Origin": []string{"http://localhost:3000"},
		})
		origin, err := GetOrigin(c)
		s.Equal("http://localhost:3000", origin)
		s.NoError(err)
	})

	s.Run("when origin has subdomain", func() {
		c := s.buildContext(url.URL{}, http.Header{
			"Origin": []string{"http://www.demoticker.org/"},
		})
		origin, err := GetOrigin(c)
		s.Equal("http://www.demoticker.org", origin)
		s.NoError(err)
	})

	s.Run("when query param is set", func() {
		c := s.buildContext(url.URL{RawQuery: "origin=http://another.demoticker.org"}, http.Header{
			"Origin": []string{"http://www.demoticker.org/"},
		})
		domain, err := GetOrigin(c)
		s.Equal("http://another.demoticker.org", domain)
		s.NoError(err)
	})
}

func (s *UtilTestSuite) TestMe() {
	s.Run("when me is not set", func() {
		c := &gin.Context{}
		_, err := Me(c)
		s.Equal("me not found", err.Error())
	})

	s.Run("when me is set", func() {
		c := &gin.Context{}
		c.Set("me", storage.User{})
		_, err := Me(c)
		s.NoError(err)
	})
}

func (s *UtilTestSuite) TestIsAdmin() {
	s.Run("when me is not set", func() {
		c := &gin.Context{}
		isAdmin := IsAdmin(c)
		s.False(isAdmin)
	})

	s.Run("when me is set", func() {
		c := &gin.Context{}
		c.Set("me", storage.User{IsSuperAdmin: true})
		isAdmin := IsAdmin(c)
		s.True(isAdmin)
	})
}

func (s *UtilTestSuite) TestTicker() {
	s.Run("when ticker is not set", func() {
		c := &gin.Context{}
		_, err := Ticker(c)
		s.Equal("ticker not found", err.Error())
	})

	s.Run("when ticker is set", func() {
		c := &gin.Context{}
		c.Set("ticker", storage.Ticker{})
		_, err := Ticker(c)
		s.NoError(err)
	})
}

func (s *UtilTestSuite) TestMessage() {
	s.Run("when message is not set", func() {
		c := &gin.Context{}
		_, err := Message(c)
		s.Equal("message not found", err.Error())
	})

	s.Run("when message is set", func() {
		c := &gin.Context{}
		c.Set("message", storage.Message{})
		_, err := Message(c)
		s.NoError(err)
	})
}

func (s *UtilTestSuite) TestUser() {
	s.Run("when user is not set", func() {
		c := &gin.Context{}
		_, err := User(c)
		s.Equal("user not found", err.Error())
	})

	s.Run("when user is set", func() {
		c := &gin.Context{}
		c.Set("user", storage.User{})
		_, err := User(c)
		s.NoError(err)
	})
}

func (s *UtilTestSuite) buildContext(u url.URL, headers http.Header) *gin.Context {
	req := http.Request{
		Header: headers,
		URL:    &u,
	}

	return &gin.Context{Request: &req}
}

func TestUtilTestSuite(t *testing.T) {
	suite.Run(t, new(UtilTestSuite))
}
