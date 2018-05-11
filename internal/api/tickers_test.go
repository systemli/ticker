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
)

func TestGetTickers(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/tickers").
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)
		assert.Equal(t, `{"data":{"tickers":[]},"status":"success","error":null}`, strings.TrimSpace(r.Body.String()))
	})
}

func TestGetTicker(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/tickers/1").
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
			"email": "admin@systemli.org"
		}
	}`

	r.POST("/v1/admin/tickers").
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
	})
}

func TestPutTicker(t *testing.T) {
	r := setup()

	var ticker model.Ticker

	ticker.Domain = "demoticker.org"

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
		SetBody(body).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
	})

	r.PUT("/v1/admin/tickers/1").
		SetBody(`malicious data`).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 400, r.Code)
	})

	r.PUT("/v1/admin/tickers/1").
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

	ticker := model.NewTicker()

	storage.DB.Save(&ticker)

	r.DELETE("/v1/admin/tickers/2").
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
	})

	r.DELETE("/v1/admin/tickers/1").
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

	return gofight.New()
}
