package util_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	. "git.codecoop.org/systemli/ticker/internal/util"
)

func TestDefaultPagination(t *testing.T) {
	req := http.Request{
		URL: &url.URL{
			RawQuery: ``,
		},
	}

	c := gin.Context{Request: &req,}
	p := NewPagination(&c)

	assert.Equal(t, p.GetLimit(), 10)
	assert.Equal(t, p.GetBefore(), 0)
	assert.Equal(t, p.GetAfter(), 0)
}

func TestCustomPagination(t *testing.T) {
	req := http.Request{
		URL: &url.URL{
			RawQuery: `limit=20&before=1&after=1`,
		},
	}

	c := gin.Context{Request: &req,}
	p := NewPagination(&c)

	assert.Equal(t, p.GetLimit(), 20)
	assert.Equal(t, p.GetBefore(), 1)
	assert.Equal(t, p.GetAfter(), 1)
}
