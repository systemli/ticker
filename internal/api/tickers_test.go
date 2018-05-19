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
)

var Token string

func TestGetTickers(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/tickers").
		SetHeader(map[string]string{"Authorization": "Bearer " + Token}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)
		assert.Equal(t, `{"data":{"tickers":[]},"status":"success","error":null}`, strings.TrimSpace(r.Body.String()))
	})
}

func TestGetTicker(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + Token}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1001,"message":"not found"}}`, strings.TrimSpace(r.Body.String()))
	})
}

func TestPostTicker(t *testing.T) {
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
		SetHeader(map[string]string{"Authorization": "Bearer " + Token}).
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
}

func TestPutTicker(t *testing.T) {
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
		SetHeader(map[string]string{"Authorization": "Bearer " + Token}).
		SetBody(body).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
	})

	r.PUT("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + Token}).
		SetBody(`malicious data`).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 400, r.Code)
	})

	r.PUT("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + Token}).
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
}

func TestDeleteTicker(t *testing.T) {
	r := setup()

	ticker := model.Ticker{
		ID:     1,
		Active: true,
	}

	storage.DB.Save(&ticker)

	r.DELETE("/v1/admin/tickers/2").
		SetHeader(map[string]string{"Authorization": "Bearer " + Token}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
	})

	r.DELETE("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + Token}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)

		type jsonResp struct {
			Data   map[string]model.Message `json:"data"`
			Status string                   `json:"status"`
			Error  interface{}              `json:"error"`
		}

		var jres jsonResp

		err := json.Unmarshal(r.Body.Bytes(), &jres)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.ResponseSuccess, jres.Status)
		assert.Nil(t, jres.Data)
		assert.Nil(t, jres.Error)
	})
}

func setup() *gofight.RequestConfig {
	gin.SetMode(gin.TestMode)

	if storage.DB == nil {
		storage.DB = storage.OpenDB("ticker_test.db")
	}
	storage.DB.Drop("Ticker")
	storage.DB.Drop("Message")

	r := gofight.New()

	if Token == "" {
		r.POST("/v1/admin/login").
			SetBody(`{"username":"admin", "password":"admin"}`).
			Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			type res struct {
				Code   int       `json:"code"`
				Expire time.Time `json:"expire"`
				Token  string    `json:"token"`
			}

			var response res
			json.Unmarshal(r.Body.Bytes(), &response)

			Token = response.Token
		})

	}

	return r
}
