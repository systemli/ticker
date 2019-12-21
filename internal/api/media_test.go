package api_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/api"
	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

func TestGetMedia(t *testing.T) {
	r := setup()

	ticker := model.Ticker{
		ID:     1,
		Active: true,
		Domain: "demoticker.org",
	}

	_ = storage.DB.Save(&ticker)

	var url string
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
			assert.Equal(t, "image/jpeg", response.Data["uploads"][0].ContentType)
			assert.NotNil(t, response.Data["uploads"][0].URL)
			assert.NotNil(t, response.Data["uploads"][0].UUID)
			assert.NotNil(t, response.Data["uploads"][0].CreationDate)
			assert.NotNil(t, response.Data["uploads"][0].ID)

			url = response.Data["uploads"][0].URL
		})

	r.GET(url).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			assert.Equal(t, "image/jpeg", r.HeaderMap.Get("Content-Type"))
			assert.Equal(t, "62497", r.HeaderMap.Get("Content-Length"))
		})

	r.GET("/media/nonexisting").
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})

	r.GET("/media/ed79e414-c399-49f8-9d49-9387df6e2768.jpg").
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})

	upload := model.NewUpload("image.jpg", "image/jpeg", 1)
	_ = storage.DB.Save(upload)

	r.GET(fmt.Sprintf("/media/%s", upload.FileName())).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})
}
