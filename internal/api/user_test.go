package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetUsersStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{IsSuperAdmin: true})
	s := &storage.MockStorage{}
	s.On("FindUsers", mock.Anything).Return([]storage.User{}, errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.GetUsers(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUsers(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{IsSuperAdmin: true})
	s := &storage.MockStorage{}
	s.On("FindUsers", mock.Anything).Return([]storage.User{}, nil)

	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.GetUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserMissingParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.GetUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUserInsufficentPermission(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1})
	c.Set("me", storage.User{ID: 2, IsSuperAdmin: false})

	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.GetUser(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetUserStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("userID", "1")
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})

	s := &storage.MockStorage{}
	s.On("FindUserByID", mock.Anything, mock.Anything).Return(storage.User{}, errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.GetUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUserMissingPermission(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 2})
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: false})

	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.GetUser(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1})
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: false})

	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.GetUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUserMissingBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{}

	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.PostUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostUserTooLongPassword(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	json := `{"email":"louis@systemli.org","password":"swusp-dud-gust-grong-yuz-swuft-plaft-glact-skast-swem-yen-kom-tut-prisp-gont"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.PostUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostUserStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	json := `{"email":"louis@systemli.org","password":"password1234"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	s := &storage.MockStorage{}
	s.On("SaveUser", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.PostUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	json := `{"email":"louis@systemli.org","password":"password1234"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	s := &storage.MockStorage{}
	s.On("SaveUser", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.PostUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutUserNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.PutUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPutUserMissingBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	c.Set("user", storage.User{})
	body := `broken_json`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.PutUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutUserStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	c.Set("user", storage.User{})
	json := `{"email":"louis@systemli.org","password":"password1234","isSuperAdmin":true,"tickers":[{"id":1}]}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("FindTickersByIDs", mock.Anything).Return([]storage.Ticker{}, nil)
	s.On("SaveUser", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.PutUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutUserStorageError2(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	c.Set("user", storage.User{})
	json := `{"email":"louis@systemli.org","password":"password1234","isSuperAdmin":true,"tickers":[{"id":1}]}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("FindTickersByIDs", mock.Anything).Return([]storage.Ticker{}, errors.New("storage error"))
	s.On("SaveUser", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.PutUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	c.Set("user", storage.User{})
	json := `{"email":"louis@systemli.org","password":"password1234","isSuperAdmin":true,"tickers":[{"id":1}]}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("FindTickersByIDs", mock.Anything).Return([]storage.Ticker{{ID: 1}}, nil)
	s.On("SaveUser", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.PutUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteUserMissingParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.DeleteUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteUserSelfUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	c.Set("user", storage.User{ID: 1})
	s := &storage.MockStorage{}
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.DeleteUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteUserStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	c.Set("user", storage.User{ID: 2})
	s := &storage.MockStorage{}
	s.On("DeleteUser", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.DeleteUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, IsSuperAdmin: true})
	c.Set("user", storage.User{ID: 2})
	s := &storage.MockStorage{}
	s.On("DeleteUser", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}

	h.DeleteUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutMeUnauthenticated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	json := `{"password":"password1234","newPassword":"password5678"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/users/me", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("SaveUser", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}
	h.PutMe(c)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPutMeFormError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, EncryptedPassword: "$2a$10$3rj/kzMI7gKPoBtJFG55tuzA.RQGYqbYQdM69LPyU.2YkGbkRu.T2"})
	json := `{"wrongparameter":"password1234","newPassword":"password5678"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/users/me", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("SaveUser", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}
	h.PutMe(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutMeWrongPassword(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, EncryptedPassword: "$2a$10$3rj/kzMI7gKPoBtJFG55tuzA.RQGYqbYQdM69LPyU.2YkGbkRu.T2"})
	json := `{"password":"wrongpassword","newPassword":"password5678"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/users/me", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("SaveUser", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}
	h.PutMe(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutMeStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, EncryptedPassword: "$2a$10$3rj/kzMI7gKPoBtJFG55tuzA.RQGYqbYQdM69LPyU.2YkGbkRu.T2"})
	json := `{"password":"password1234","newPassword":"password5678"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/users/me", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("SaveUser", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}
	h.PutMe(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutMeOk(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("me", storage.User{ID: 1, EncryptedPassword: "$2a$10$3rj/kzMI7gKPoBtJFG55tuzA.RQGYqbYQdM69LPyU.2YkGbkRu.T2"})
	json := `{"password":"password1234","newPassword":"password5678"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/admin/users/me", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")
	s := &storage.MockStorage{}
	s.On("SaveUser", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.LoadConfig(""),
	}
	h.PutMe(c)
	assert.Equal(t, http.StatusOK, w.Code)
}
