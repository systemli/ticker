package helper

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/systemli/ticker/internal/storage"
)

func TestGetDomainEmptyOrigin(t *testing.T) {
	req := http.Request{
		URL: &url.URL{},
	}

	c := gin.Context{Request: &req}

	domain, err := GetDomain(&c)
	assert.Equal(t, "", domain)
	assert.Equal(t, "Origin header not found", err.Error())
}

func TestGetDomainLocalhost(t *testing.T) {
	req := http.Request{
		Header: http.Header{
			"Origin": []string{"http://localhost/"},
		},
		URL: &url.URL{},
	}

	c := gin.Context{Request: &req}

	domain, err := GetDomain(&c)
	assert.Equal(t, "localhost", domain)
	assert.Equal(t, nil, err)
}

func TestGetDomainLocalhostPort(t *testing.T) {
	req := http.Request{
		Header: http.Header{
			"Origin": []string{"http://localhost:3000/"},
		},
		URL: &url.URL{},
	}

	c := gin.Context{Request: &req}

	domain, err := GetDomain(&c)
	assert.Equal(t, "localhost", domain)
	assert.Equal(t, nil, err)
}

func TestGetDomainWWW(t *testing.T) {
	req := http.Request{
		Header: http.Header{
			"Origin": []string{"http://www.demoticker.org/"},
		},
		URL: &url.URL{},
	}

	c := gin.Context{Request: &req}

	domain, err := GetDomain(&c)
	assert.Equal(t, "demoticker.org", domain)
	assert.Equal(t, nil, err)
}

func TestGetDomainOriginQueryOverwrite(t *testing.T) {
	req := http.Request{
		Header: http.Header{
			"Origin": []string{"http://www.demoticker.org/"},
		},
		URL: &url.URL{RawQuery: "origin=another.demoticker.org"},
	}

	c := gin.Context{Request: &req}

	domain, err := GetDomain(&c)
	assert.Equal(t, "another.demoticker.org", domain)
	assert.Equal(t, nil, err)
}

func TestMe(t *testing.T) {
	c := &gin.Context{}
	_, err := Me(c)

	assert.NotNil(t, err)

	c.Set("user", storage.User{})

	_, err = Me(c)

	assert.Nil(t, err)
}

func TestIsAdmin(t *testing.T) {
	c := &gin.Context{}
	isAdmin := IsAdmin(c)

	assert.False(t, isAdmin)

	c.Set("user", storage.User{IsSuperAdmin: true})

	isAdmin = IsAdmin(c)

	assert.True(t, isAdmin)
}

func TestTicker(t *testing.T) {
	c := &gin.Context{}

	_, err := Ticker(c)
	assert.NotNil(t, err)

	c.Set("ticker", storage.Ticker{})

	_, err = Ticker(c)
	assert.Nil(t, err)
}
