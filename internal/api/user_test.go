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

func TestGetUsersForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetUsers(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetUsersStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	s.On("FindUsers").Return([]storage.User{}, errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetUsers(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUsers(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	s.On("FindUsers").Return([]storage.User{}, nil)

	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserMissingParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserMissingUserInContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("userID", "1")
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserInsufficentPermission(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("userID", "1")
	c.Set("user", storage.User{ID: 2, IsSuperAdmin: false})

	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetUser(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetUserStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("userID", "1")
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})

	s := &storage.MockTickerStorage{}
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.AddParam("userID", "1")
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})

	s := &storage.MockTickerStorage{}
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.GetUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUserAsNonAdmin(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: false})

	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostUser(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPostUserMissingBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	c.Request = &http.Request{}

	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
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
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	s.On("SaveUser", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
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
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	s.On("SaveUser", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PostUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutUserAsNonAdmin(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: false})

	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutUser(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPutUserMissingParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutUserNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	c.AddParam("userID", "2")
	s := &storage.MockTickerStorage{}
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, errors.New("not found"))

	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPutUserMissingBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	c.AddParam("userID", "2")
	body := `broken_json`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", strings.NewReader(body))
	c.Request.Header.Add("Content-Type", "application/json")

	s := &storage.MockTickerStorage{}
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, nil)

	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutUserStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	c.AddParam("userID", "2")
	json := `{"email":"louis@systemli.org","password":"password1234","is_super_admin":true,"tickers":[1]}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")

	s := &storage.MockTickerStorage{}
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, nil)
	s.On("SaveUser", mock.Anything).Return(errors.New("storage error"))

	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	c.AddParam("userID", "2")
	json := `{"email":"louis@systemli.org","password":"password1234","is_super_admin":true,"tickers":[1]}`
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/admin/users", strings.NewReader(json))
	c.Request.Header.Add("Content-Type", "application/json")

	s := &storage.MockTickerStorage{}
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, nil)
	s.On("SaveUser", mock.Anything).Return(nil)

	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.PutUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteUserAsNonAdmin(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: false})

	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteUser(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteUserMissingParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteUserSelfUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	c.AddParam("userID", "1")

	s := &storage.MockTickerStorage{}
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteUserNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	c.AddParam("userID", "2")

	s := &storage.MockTickerStorage{}
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, errors.New("not found"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteUserStorageError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	c.AddParam("userID", "2")

	s := &storage.MockTickerStorage{}
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, nil)
	s.On("DeleteUser", mock.Anything).Return(errors.New("storage error"))
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", storage.User{ID: 1, IsSuperAdmin: true})
	c.AddParam("userID", "2")

	s := &storage.MockTickerStorage{}
	s.On("FindUserByID", mock.Anything).Return(storage.User{}, nil)
	s.On("DeleteUser", mock.Anything).Return(nil)
	h := handler{
		storage: s,
		config:  config.NewConfig(),
	}

	h.DeleteUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
