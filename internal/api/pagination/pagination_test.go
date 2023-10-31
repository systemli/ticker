package pagination

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type PaginationTestSuite struct {
	suite.Suite
}

func (s *PaginationTestSuite) TestNewPagination() {
	s.Run("with default values", func() {
		req := http.Request{
			URL: &url.URL{
				RawQuery: ``,
			},
		}

		c := gin.Context{Request: &req}
		p := NewPagination(&c)

		s.Equal(10, p.GetLimit())
		s.Equal(0, p.GetBefore())
		s.Equal(0, p.GetAfter())
	})

	s.Run("with custom values", func() {
		req := http.Request{
			URL: &url.URL{
				RawQuery: `limit=20&before=1&after=1`,
			},
		}

		c := gin.Context{Request: &req}
		p := NewPagination(&c)

		s.Equal(20, p.GetLimit())
		s.Equal(1, p.GetBefore())
		s.Equal(1, p.GetAfter())
	})
}

func TestPaginationTestSuite(t *testing.T) {
	suite.Run(t, new(PaginationTestSuite))
}
