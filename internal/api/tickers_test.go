package api_test

import (
	"testing"
	"encoding/json"

	"github.com/appleboy/gofight"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"git.codecoop.org/systemli/ticker/internal/api"
	"git.codecoop.org/systemli/ticker/internal/model"
	"git.codecoop.org/systemli/ticker/internal/storage"
	"strings"
	"time"
	"fmt"
)

var AdminToken string
var UserToken string

func TestGetTickersHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/tickers").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)
		assert.Equal(t, `{"data":{"tickers":null},"status":"success","error":null}`, strings.TrimSpace(r.Body.String()))
	})

	r.GET("/v1/admin/tickers").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)
		assert.Equal(t, `{"data":{"tickers":null},"status":"success","error":null}`, strings.TrimSpace(r.Body.String()))
	})
}

func TestGetTickerHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1001,"message":"not found"}}`, strings.TrimSpace(r.Body.String()))
	})

	r.GET("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 403, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
	})
}

func TestPostTickerHandler(t *testing.T) {
	r := setup()

	body := `{
		"title": "Ticker",
		"domain": "prozessticker.org",
		"description": "Beschreibung",
		"active": true,
		"information": {
			"url": "https://www.systemli.org",
			"email": "admin@systemli.org",
			"twitter": "systemli"
		}
	}`

	r.POST("/v1/admin/tickers").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)

		type jsonResp struct {
			Data   map[string]model.Ticker `json:"data"`
			Status string                  `json:"status"`
			Error  interface{}             `json:"error"`
		}

		var jres jsonResp

		err := json.Unmarshal(r.Body.Bytes(), &jres)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.ResponseSuccess, jres.Status)
		assert.Equal(t, nil, jres.Error)
		assert.Equal(t, 1, len(jres.Data))

		ticker := jres.Data["ticker"]

		assert.Equal(t, "Ticker", ticker.Title)
		assert.Equal(t, "prozessticker.org", ticker.Domain)
		assert.Equal(t, true, ticker.Active)
		assert.Equal(t, "https://www.systemli.org", ticker.Information.URL)
		assert.Equal(t, "admin@systemli.org", ticker.Information.Email)
		assert.Equal(t, "systemli", ticker.Information.Twitter)
	})

	r.POST("/v1/admin/tickers").
		SetBody(body).
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 403, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
	})
}

func TestPutTickerHandler(t *testing.T) {
	r := setup()

	ticker := model.Ticker{
		ID:     1,
		Active: true,
		Domain: "demoticker.org",
	}

	storage.DB.Save(&ticker)

	body := `{
		"title": "Ticker",
		"domain": "prozessticker.org",
		"description": "Beschreibung",
		"active": false,
		"information": {
			"url": "https://www.systemli.org",
			"email": "admin@systemli.org"
		}
	}`

	r.PUT("/v1/admin/tickers/100").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
	})

	r.PUT("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(`malicious data`).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 400, r.Code)
	})

	r.PUT("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)

		type jsonResp struct {
			Data   map[string]model.Ticker `json:"data"`
			Status string                  `json:"status"`
			Error  interface{}             `json:"error"`
		}

		var jres jsonResp

		err := json.Unmarshal(r.Body.Bytes(), &jres)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.ResponseSuccess, jres.Status)
		assert.Equal(t, nil, jres.Error)
		assert.Equal(t, 1, len(jres.Data))

		ticker := jres.Data["ticker"]

		assert.Equal(t, 1, ticker.ID)
		assert.Equal(t, "Ticker", ticker.Title)
		assert.Equal(t, "prozessticker.org", ticker.Domain)
		assert.Equal(t, false, ticker.Active)
	})

	r.PUT("/v1/admin/tickers/1").
		SetBody(body).
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 403, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
	})
}

func TestDeleteTickerHandler(t *testing.T) {
	r := setup()

	ticker := model.Ticker{
		ID:     1,
		Active: true,
	}

	storage.DB.Save(&ticker)

	r.DELETE("/v1/admin/tickers/2").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
	})

	r.DELETE("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)

		var jres struct {
			Data   map[string]model.Message `json:"data"`
			Status string                   `json:"status"`
			Error  interface{}              `json:"error"`
		}

		err := json.Unmarshal(r.Body.Bytes(), &jres)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.ResponseSuccess, jres.Status)
		assert.Nil(t, jres.Data)
		assert.Nil(t, jres.Error)
	})

	r.DELETE("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 403, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
	})
}

func setup() *gofight.RequestConfig {
	gin.SetMode(gin.TestMode)

	model.Config = model.NewConfig()

	if storage.DB == nil {
		storage.DB = storage.OpenDB("ticker_test.db")
	}
	storage.DB.Drop("Ticker")
	storage.DB.Drop("Message")
	storage.DB.Drop("User")

	admin, _ := model.NewUser("admin@systemli.org", "password")
	admin.IsSuperAdmin = true

	storage.DB.Save(admin)

	user, _ := model.NewUser("louis@systemli.org", "password")
	storage.DB.Save(user)

	if AdminToken == "" {
		AdminToken = token("admin@systemli.org", "password")
	}
	if UserToken == "" {
		UserToken = token("louis@systemli.org", "password")
	}

	return gofight.New()
}

func token(username, password string) string {
	var token string

	r := gofight.New()
	r.POST("/v1/admin/login").
		SetBody(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password)).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

		var response struct {
			Code   int       `json:"code"`
			Expire time.Time `json:"expire"`
			Token  string    `json:"token"`
		}

		json.Unmarshal(r.Body.Bytes(), &response)

		token = response.Token
	})

	return token
}