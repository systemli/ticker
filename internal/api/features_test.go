package api

import (
	"strings"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetFeaturesHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/features").
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 401, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1002,"message":"auth header is empty"}}`, strings.TrimSpace(r.Body.String()))
		})

	r.GET("/v1/admin/features").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
		})

	r.GET("/v1/admin/features").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			assert.Equal(t, `{"data":{"features":{"telegram_enabled":false,"twitter_enabled":false}},"status":"success","error":null}`, strings.TrimSpace(r.Body.String()))
		})
}
