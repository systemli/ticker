package api_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/api"
	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

type uploadResponse struct {
	Data   map[string][]model.UploadResponse `json:"data"`
	Status string                            `json:"status"`
	Error  map[string]interface{}            `json:"error"`
}

func TestPostUploadSuccessful(t *testing.T) {
	r := setup()

	ticker := initUploadTestData()

	r.POST("/v1/admin/upload").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetFileFromPath([]gofight.UploadFile{{Name: "files", Path: "../../testdata/gopher.jpg"}}, gofight.H{"ticker": strconv.Itoa(ticker.ID)}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			var response uploadResponse
			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, response.Status)
			assert.Equal(t, 1, len(response.Data))
			assert.Equal(t, 1, len(response.Data["uploads"]))
			assert.NotNil(t, response.Data["uploads"][0].UUID)
			assert.NotNil(t, response.Data["uploads"][0].ID)
			assert.NotNil(t, response.Data["uploads"][0].URL)
			assert.NotNil(t, response.Data["uploads"][0].CreationDate)
		})
}

func TestPostUploadGIF(t *testing.T) {
	r := setup()

	ticker := initUploadTestData()
	files := []gofight.UploadFile{{Name: "files", Path: "../../testdata/gopher-dance.gif"}}

	r.POST("/v1/admin/upload").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetFileFromPath(files, gofight.H{"ticker": strconv.Itoa(ticker.ID)}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
		})
}

func TestPostUploadTickerNonExisting(t *testing.T) {
	r := setup()

	r.POST("/v1/admin/upload").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetFileFromPath([]gofight.UploadFile{{Name: "files", Path: "../../testdata/gopher.jpg"}}, gofight.H{"ticker": "2"}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)

			var response uploadResponse

			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseError, response.Status)
		})
}

func TestPostUploadUnauthorized(t *testing.T) {
	r := setup()

	ticker := initUploadTestData()

	r.POST("/v1/admin/upload").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		SetFileFromPath([]gofight.UploadFile{{Name: "files", Path: "../../testdata/gopher.jpg"}}, gofight.H{"ticker": strconv.Itoa(ticker.ID)}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)

			var response uploadResponse

			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseError, response.Status)
		})
}

func TestPostUploadWrongContentType(t *testing.T) {
	r := setup()

	ticker := initUploadTestData()
	r.POST("/v1/admin/upload").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetFileFromPath([]gofight.UploadFile{{Name: "files", Path: "../../README.md"}}, gofight.H{"ticker": strconv.Itoa(ticker.ID)}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)

			var response uploadResponse

			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseError, response.Status)
		})
}

func TestPostUploadMissingTicker(t *testing.T) {
	r := setup()

	r.POST("/v1/admin/upload").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetFileFromPath([]gofight.UploadFile{{Name: "files", Path: "../../README.md"}}, gofight.H{}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func TestPostUploadWrongTickerParam(t *testing.T) {
	r := setup()

	r.POST("/v1/admin/upload").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetFileFromPath([]gofight.UploadFile{{Name: "files", Path: "../../README.md"}}, gofight.H{"ticker": "string"}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func TestPostUploadMissingFiles(t *testing.T) {
	r := setup()

	ticker := initUploadTestData()

	r.POST("/v1/admin/upload").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetFileFromPath([]gofight.UploadFile{}, gofight.H{"ticker": strconv.Itoa(ticker.ID)}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func TestPostUploadTooMuchFiles(t *testing.T) {
	r := setup()

	ticker := initUploadTestData()

	files := []gofight.UploadFile{
		{Name: "files", Path: "../../testdata/gopher.jpg"},
		{Name: "files", Path: "../../testdata/gopher.jpg"},
		{Name: "files", Path: "../../testdata/gopher.jpg"},
		{Name: "files", Path: "../../testdata/gopher.jpg"},
	}

	r.POST("/v1/admin/upload").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetFileFromPath(files, gofight.H{"ticker": strconv.Itoa(ticker.ID)}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)

			var response uploadResponse

			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseError, response.Status)
			assert.Equal(t, model.ErrorTooMuchFiles, response.Error["message"])
		})
}

func initUploadTestData() *model.Ticker {
	ticker := &model.Ticker{
		ID:     1,
		Active: true,
		Domain: "demoticker.org",
	}

	_ = storage.DB.Save(ticker)

	return ticker
}
