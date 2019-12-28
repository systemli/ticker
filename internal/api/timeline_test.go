package api_test

import (
	"encoding/json"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/api"
	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

type timelineResponse struct {
	Data   map[string][]model.MessageResponse `json:"data"`
	Status string                             `json:"status"`
	Error  interface{}                        `json:"error"`
}

func TestGetTimelineHandler(t *testing.T) {
	r := setup()
	ticker := model.NewTicker()
	ticker.Active = true
	ticker.Domain = "localhost"

	_ = storage.DB.Save(ticker)

	r.GET("/v1/timeline?origin="+ticker.Domain).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			var response timelineResponse
			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, "success", response.Status)
			assert.Equal(t, map[string][]model.MessageResponse{"messages": []model.MessageResponse(nil)}, response.Data)
			assert.Nil(t, response.Error)
		})
}

func TestGetTimelineHandler2(t *testing.T) {
	r := setup()

	r.GET("/v1/timeline?origin=non").
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			var response timelineResponse
			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, "error", response.Status)
			assert.Equal(t, map[string][]model.MessageResponse{"messages": []model.MessageResponse(nil)}, response.Data)
			assert.NotNil(t, response.Error)
		})
}
