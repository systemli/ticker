package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/storage"
)

type AdminTestSuite struct {
	suite.Suite
}

func (s *AdminTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *AdminTestSuite) TestNeedAdmin() {
	s.Run("when user is missing", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		mw := NeedAdmin()

		mw(c)

		s.Equal(http.StatusBadRequest, w.Code)
	})

	s.Run("when user is not admin", func() {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("me", storage.User{})
		mw := NeedAdmin()

		mw(c)

		s.Equal(http.StatusForbidden, w.Code)
	})
}

func TestAdminTestSuite(t *testing.T) {
	suite.Run(t, new(AdminTestSuite))
}
