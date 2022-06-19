package api

import (
	"encoding/json"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

func TestGetTimelineHandler(t *testing.T) {
	r := setup()
	ticker := model.NewTicker()
	ticker.Active = true
	ticker.Domain = "localhost"

	_ = storage.DB.Save(ticker)

	r.GET("/v1/timeline?origin="+ticker.Domain).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			var response timelineResponse
			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, "success", response.Status)
			assert.Equal(t, []message(nil), response.Data.Messages)
		})
}

func TestGetTimelineHandler2(t *testing.T) {
	r := setup()

	r.GET("/v1/timeline?origin=non").
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)

			var response errorResponse
			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, "error", response.Status)
			assert.Equal(t, "Could not find a ticker.", response.Error.Message)
			assert.Equal(t, 1000, response.Error.Code)
		})
}
