package api_test

import (
	"testing"
	"net/http"
	"git.codecoop.org/systemli/ticker/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/url"
)

func TestGetDomainEmptyOrigin(t *testing.T) {
	req := http.Request{}
	req.URL = &url.URL{}

	c := gin.Context{Request: &req,}

	domain, err := api.GetDomain(&c)
	assert.Equal(t, "", domain)
	assert.Equal(t, "Origin header not found", err.Error())
}

func TestGetDomainLocalhost(t *testing.T) {
	req := http.Request{
		Header: http.Header{
			"Origin": []string{"http://localhost/"},
		},
	}

	c := gin.Context{Request: &req,}

	domain, err := api.GetDomain(&c)
	assert.Equal(t, "localhost", domain)
	assert.Equal(t, nil, err)
}

func TestGetDomainLocalhostPort(t *testing.T) {
	req := http.Request{
		Header: http.Header{
			"Origin": []string{"http://localhost:3000/"},
		},
	}

	c := gin.Context{Request: &req,}

	domain, err := api.GetDomain(&c)
	assert.Equal(t, "localhost", domain)
	assert.Equal(t, nil, err)
}

func TestGetDomainWWW(t *testing.T) {
	req := http.Request{
		Header: http.Header{
			"Origin": []string{"http://www.demoticker.org/"},
		},
	}

	c := gin.Context{Request: &req,}

	domain, err := api.GetDomain(&c)
	assert.Equal(t, "demoticker.org", domain)
	assert.Equal(t, nil, err)
}
