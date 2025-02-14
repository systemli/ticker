package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/url"
	"testing"
)

type PrometheusTestSuite struct {
	suite.Suite
}

func (s *PrometheusTestSuite) TestPrepareOrigin() {

	s.Run("when request is empty", func() {
		origin := prepareOrigin(s.buildContext(url.URL{}, http.Header{}))
		s.Empty(origin)
	})

	s.Run("when origin is in query", func() {
		s.Run("when origin is valid url", func() {
			origin := prepareOrigin(s.buildContext(url.URL{RawQuery: "origin=https://example.com"}, http.Header{}))
			s.Equal("example.com", origin)
		})

		s.Run("when origin is invalid url", func() {
			origin := prepareOrigin(s.buildContext(url.URL{RawQuery: "origin=invalid"}, http.Header{}))
			s.Empty(origin)
		})
	})

	s.Run("when origin is in header", func() {
		s.Run("when origin is valid url", func() {
			origin := prepareOrigin(s.buildContext(url.URL{}, http.Header{
				"Origin": []string{"https://example.com"},
			}))
			s.Equal("example.com", origin)
		})

		s.Run("when origin is invalid url", func() {
			origin := prepareOrigin(s.buildContext(url.URL{}, http.Header{
				"Origin": []string{"invalid"},
			}))
			s.Empty(origin)
		})
	})
}

func (s *PrometheusTestSuite) buildContext(u url.URL, headers http.Header) *gin.Context {
	req := http.Request{
		Header: headers,
		URL:    &u,
	}

	return &gin.Context{Request: &req}
}

func TestPrometheusSuite(t *testing.T) {
	suite.Run(t, new(PrometheusTestSuite))
}
