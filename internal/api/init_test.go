package api_test

import (
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"

	"encoding/json"
	"github.com/systemli/ticker/internal/api"
	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

func TestGetInitHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/init").
		SetHeader(map[string]string{"Origin": "http://www.demoticker.org/"}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, r.Code, 200)

			type res struct {
				Data   map[string]interface{} `json:"data"`
				Status string                 `json:"status"`
				Error  interface{}            `json:"error"`
			}

			var response res

			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, response.Status)
			assert.Equal(t, nil, response.Error)
			assert.Equal(t, 2, len(response.Data))

			settings := response.Data["settings"].(map[string]interface{})
			assert.Equal(t, float64(10000), settings["refresh_interval"])

			ticker := response.Data["ticker"]
			assert.Nil(t, ticker)
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
		Run(api.API(), func(response gofight.HTTPResponse, request gofight.HTTPRequest) {
			assert.Equal(t, response.Code, 200)

			var data struct {
				Data   map[string]interface{}
				Status string
				Error  map[string]interface{}
			}

			err := json.Unmarshal(response.Body.Bytes(), &data)
			if err != nil {
				t.Fail()
			}

			assert.Nil(t, data.Error)
			assert.Equal(t, model.ResponseSuccess, data.Status)
			assert.Equal(t, 2, len(data.Data))

			ticker := data.Data["ticker"]
			assert.NotNil(t, ticker)
		})
}
