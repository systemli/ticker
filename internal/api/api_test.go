package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func init() {
	logrus.SetOutput(io.Discard)
	gin.SetMode(gin.TestMode)
}

func TestHealthz(t *testing.T) {
	c := config.NewConfig()
	s := &storage.MockTickerStorage{}
	r := API(c, s)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginNotSuccessful(t *testing.T) {
	c := config.NewConfig()
	s := &storage.MockTickerStorage{}
	user := storage.User{}
	user.UpdatePassword("password")
	s.On("FindUserByEmail", mock.Anything).Return(user, nil)
	r := API(c, s)

	body := `{"username":"louis@systemli.org","password":"WRONG"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/admin/login", strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var res response.Response
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Nil(t, res.Data)
	assert.Equal(t, res.Error.Code, response.CodeBadCredentials)
	assert.Equal(t, res.Error.Message, response.Unauthorized)
}

func TestLoginSuccessful(t *testing.T) {
	c := config.NewConfig()
	s := &storage.MockTickerStorage{}
	user := storage.User{}
	user.UpdatePassword("password")
	s.On("FindUserByEmail", mock.Anything).Return(user, nil)
	r := API(c, s)

	body := `{"username":"louis@systemli.org","password":"password"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/admin/login", strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var res struct {
		Code   int       `json:"code"`
		Expire time.Time `json:"expire"`
		Token  string    `json:"token"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.NotNil(t, res.Expire)
	assert.NotNil(t, res.Token)
}
