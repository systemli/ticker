package user

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/systemli/ticker/internal/storage"
)

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Middleware Suite")
}

var _ = Describe("User Middleware", func() {
	When("id is not present", func() {
		It("should return an error", func() {
			mockStorage := &storage.MockTickerStorage{}
			mw := UserMiddleware(mockStorage)
			gin.SetMode(gin.ReleaseMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			mw(c)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	When("id is present", func() {
		Context("user not found", func() {
			mockStorage := &storage.MockTickerStorage{}
			mockStorage.On("FindUserByID", mock.Anything).Return(storage.User{}, errors.New("not found"))
			mw := UserMiddleware(mockStorage)
			gin.SetMode(gin.ReleaseMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("id", float64(1))

			mw(c)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		Context("user found", func() {
			mockStorage := &storage.MockTickerStorage{}
			mockStorage.On("FindUserByID", mock.Anything).Return(storage.User{ID: 1}, nil)
			mw := UserMiddleware(mockStorage)
			gin.SetMode(gin.ReleaseMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("id", float64(1))

			mw(c)

			user, exists := c.Get("user")
			Expect(exists).To(BeTrue())
			Expect(user).To(BeAssignableToTypeOf(storage.User{}))
		})
	})
})
