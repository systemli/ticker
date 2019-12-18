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

func TestPostUpload(t *testing.T) {
	r := setup()

	ticker := model.Ticker{
		ID:     1,
		Active: true,
		Domain: "demoticker.org",
	}

	storage.DB.Save(&ticker)

	r.POST("/v1/admin/upload").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetFileFromPath([]gofight.UploadFile{{Name: "files", Path: "../../testdata/gopher.jpg"}}, gofight.H{"ticker": "1"}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			var response struct {
				Data   map[string][]model.UploadResponse `json:"data"`
				Status string                            `json:"status"`
				Error  interface{}                       `json:"error"`
			}

			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, response.Status)
			assert.Equal(t, nil, response.Error)
			assert.Equal(t, 1, len(response.Data))
			assert.Equal(t, 1, len(response.Data["uploads"]))
			assert.NotNil(t, response.Data["uploads"][0].UUID)
			assert.NotNil(t, response.Data["uploads"][0].ID)
			assert.NotNil(t, response.Data["uploads"][0].URL)
			assert.NotNil(t, response.Data["uploads"][0].CreationDate)
		})

	r.POST("/v1/admin/upload").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetFileFromPath([]gofight.UploadFile{{Name: "files", Path: "../../testdata/gopher.jpg"}}, gofight.H{"ticker": "2"}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)

			var response struct {
				Data   map[string][]model.UploadResponse `json:"data"`
				Status string                            `json:"status"`
				Error  interface{}                       `json:"error"`
			}

			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseError, response.Status)
		})

	r.POST("/v1/admin/upload").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		SetFileFromPath([]gofight.UploadFile{{Name: "files", Path: "../../testdata/gopher.jpg"}}, gofight.H{"ticker": "1"}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)

			var response struct {
				Data   map[string][]model.UploadResponse `json:"data"`
				Status string                            `json:"status"`
				Error  interface{}                       `json:"error"`
			}

			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseError, response.Status)
		})
}
