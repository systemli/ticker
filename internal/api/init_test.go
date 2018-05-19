package api_test

import (
	"testing"
	"git.codecoop.org/systemli/ticker/internal/api"
	"github.com/appleboy/gofight"
	"github.com/docker/docker/pkg/testutil/assert"
	"git.codecoop.org/systemli/ticker/internal/model"
	"git.codecoop.org/systemli/ticker/internal/storage"
	"strings"
)

func TestGetInit(t *testing.T) {
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
		assert.Equal(t, strings.TrimSpace(response.Body.String()), `{"data":{"settings":{"refresh_interval":10},"ticker":{"id":1,"creation_date":"0001-01-01T00:00:00Z","domain":"demoticker.org","title":"Demoticker","description":"Description","active":true,"information":{"author":"","url":"","email":"","twitter":"","facebook":""}}},"status":"success","error":null}`)

	})
}
