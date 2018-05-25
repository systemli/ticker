package api_test

import (
	"strings"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"

	"git.codecoop.org/systemli/ticker/internal/api"
	"git.codecoop.org/systemli/ticker/internal/model"
	"git.codecoop.org/systemli/ticker/internal/storage"
	"encoding/json"
)

func TestGetInitHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/init").
		SetHeader(map[string]string{"Origin": "http://www.demoticker.org/"}).
		Run(api.API(), func(response gofight.HTTPResponse, request gofight.HTTPRequest) {
		assert.Equal(t, response.Code, 200)
		assert.Equal(t, strings.TrimSpace(response.Body.String()), `{"data":{"settings":{"refresh_interval":10},"ticker":null},"status":"success","error":null}`)
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
		assert.Equal(t,2, len(data.Data))
	})
}
