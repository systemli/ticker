package api

import (
	"encoding/json"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

func TestGetInitHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/init").
		SetHeader(map[string]string{"Origin": "http://www.demoticker.org/"}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, 200)

			var response initResponse
			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, response.Status)
			assert.Nil(t, response.Error)
			assert.Nil(t, response.Data.Ticker)

			assert.Equal(t, 10000, response.Data.Settings.RefreshInterval)
		})

	ticker := new(model.Ticker)
	ticker.ID = 1
	ticker.Active = true
	ticker.Title = "Demoticker"
	ticker.Description = "Description"
	ticker.Domain = "demoticker.org"

	storage.DB.Save(ticker)

	r.GET("/v1/init").
		SetHeader(map[string]string{"Origin": "http://www.demoticker.org/"}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, 200)

			var response initResponse
			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Nil(t, response.Error)
			assert.Equal(t, model.ResponseSuccess, response.Status)

			assert.NotNil(t, response.Data.Ticker)
		})
}
