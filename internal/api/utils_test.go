package api_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/api"
)

func TestGetDomainEmptyOrigin(t *testing.T) {
	req := http.Request{
		URL: &url.URL{},
	}

	c := gin.Context{Request: &req}

	domain, err := api.GetDomain(&c)
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

	domain, err := api.GetDomain(&c)
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

	domain, err := api.GetDomain(&c)
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

	domain, err := api.GetDomain(&c)
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

	domain, err := api.GetDomain(&c)
	assert.Equal(t, "another.demoticker.org", domain)
	assert.Equal(t, nil, err)
}
