package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Middleware Suite")
}

var _ = Describe("Auth Middleware", func() {
	gin.SetMode(gin.TestMode)

	When("Authenticator", func() {
		Context("empty form is sent", func() {
			mockStorage := &storage.MockTickerStorage{}

			authenticator := Authenticator(mockStorage)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{}`))

			It("should return an error", func() {
				_, err := authenticator(c)
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(errors.New("missing Username or Password")))
			})
		})

		Context("user not found", func() {
			mockStorage := &storage.MockTickerStorage{}
			mockStorage.On("FindUserByEmail", mock.Anything).Return(storage.User{}, errors.New("not found"))

			authenticator := Authenticator(mockStorage)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username": "user@systemli.org", "password": "password"}`))
			c.Request.Header.Set("Content-Type", "application/json")

			It("should return an error", func() {
				_, err := authenticator(c)
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(errors.New("not found")))
			})
		})

		Context("user is found", func() {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
			user := storage.User{
				ID:                1,
				Email:             "user@systemli.org",
				EncryptedPassword: string(hashedPassword),
				IsSuperAdmin:      false,
			}
			mockStorage := &storage.MockTickerStorage{}
			mockStorage.On("FindUserByEmail", mock.Anything).Return(user, nil)

			authenticator := Authenticator(mockStorage)

			It("should return a user with correct password", func() {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username": "user@systemli.org", "password": "password"}`))
				c.Request.Header.Set("Content-Type", "application/json")
				user, err := authenticator(c)
				Expect(err).To(BeNil())
				Expect(user).To(BeAssignableToTypeOf(storage.User{}))
			})

			It("should return an error with incorrect password", func() {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username": "user@systemli.org", "password": "password1"}`))
				c.Request.Header.Set("Content-Type", "application/json")

				_, err := authenticator(c)
				Expect(err).NotTo(BeNil())
			})
		})
	})

	When("Authorizator", func() {
		Context("storage returns no user", func() {
			mockStorage := &storage.MockTickerStorage{}
			mockStorage.On("FindUserByID", mock.Anything).Return(storage.User{}, errors.New("user not found"))
			authorizator := Authorizator(mockStorage)
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			It("should return false", func() {
				found := authorizator(float64(1), c)
				Expect(found).To(BeFalse())
			})
		})

		Context("storage returns a user", func() {
			mockStorage := &storage.MockTickerStorage{}
			mockStorage.On("FindUserByID", mock.Anything).Return(storage.User{ID: 1}, nil)
			authorizator := Authorizator(mockStorage)
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			It("should return true", func() {
				found := authorizator(float64(1), c)
				Expect(found).To(BeTrue())
			})
		})
	})

	When("Unauthorized", func() {
		It("should return a 403 with json payload", func() {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Request = httptest.NewRequest(http.MethodGet, "/login", nil)

			Unauthorized(c, 403, "unauthorized")

			err := json.Unmarshal(rr.Body.Bytes(), &response.Response{})
			Expect(rr.Code).To(Equal(403))
			Expect(err).To(BeNil())
		})
	})

	When("FillClaims", func() {
		Context("invalid user is given", func() {
			It("should return empty claims", func() {
				claims := FillClaim("empty")
				Expect(claims).To(Equal(jwt.MapClaims{}))
			})
		})

		Context("valid user is given", func() {
			It("should return the claims", func() {
				user := storage.User{ID: 1, Email: "user@systemli.org", IsSuperAdmin: true}
				claims := FillClaim(user)

				Expect(claims).To(Equal(jwt.MapClaims{"id": 1, "email": "user@systemli.org", "roles": []string{"user", "admin"}}))
			})
		})
	})
})
