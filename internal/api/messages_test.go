package api

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/appleboy/gofight/v2"
	"github.com/google/uuid"
	geojson "github.com/paulmach/go.geojson"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

func TestGetMessagesHandler(t *testing.T) {
	r := setup()

	ticker := model.Ticker{
		ID:     1,
		Active: true,
	}

	storage.DB.Save(&ticker)

	r.GET("/v1/admin/tickers/1/messages").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			assert.Equal(t, `{"data":{"messages":[]},"status":"success","error":null}`, strings.TrimSpace(r.Body.String()))
		})

	r.GET("/v1/admin/tickers/1/messages").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
		})
}

func TestGetMessageHandler(t *testing.T) {
	r := setup()

	ticker := model.Ticker{
		ID:     1,
		Active: true,
	}

	storage.DB.Save(&ticker)

	r.GET("/v1/admin/tickers/1/messages/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1001,"message":"not found"}}`, strings.TrimSpace(r.Body.String()))
		})

	r.GET("/v1/admin/tickers/1/messages/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
		})

	message := model.NewMessage()
	message.Text = "text"
	message.Ticker = ticker.ID

	storage.DB.Save(message)

	r.GET("/v1/admin/tickers/1/messages/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			var response struct {
				Data   map[string]model.Message `json:"data"`
				Status string                   `json:"status"`
				Error  interface{}              `json:"error"`
			}
			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Nil(t, response.Error)
			assert.Equal(t, 1, len(response.Data))
			assert.Equal(t, "text", response.Data["message"].Text)
		})

	r.GET("/v1/admin/tickers/1/messages/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
		})

	var user model.User
	err := storage.DB.One("ID", 2, &user)
	if err != nil {
		t.Fail()
	}

	user.Tickers = []int{1}
	err = storage.DB.Save(&user)
	if err != nil {
		t.Fail()
	}

	r.GET("/v1/admin/tickers/1/messages/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			assert.Equal(t, 200, r.Code)

			var response struct {
				Data   map[string]model.Message `json:"data"`
				Status string                   `json:"status"`
				Error  interface{}              `json:"error"`
			}
			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Nil(t, response.Error)
			assert.Equal(t, 1, len(response.Data))
			assert.Equal(t, "text", response.Data["message"].Text)
		})
}

func TestPostMessageHandler(t *testing.T) {
	r := setup()

	ticker := model.Ticker{
		ID:       1,
		Active:   true,
		Hashtags: []string{`#hashtag`},
	}

	storage.DB.Save(&ticker)

	body := `{
		"text": "message",
		"geo_information": {
			"type" : "FeatureCollection",
			"features" : [{ 
				"type" : "Feature", 
				"properties" : {  
					"capacity" : "10", 
					"type" : "U-Rack",
					"mount" : "Surface"
				}, 
				"geometry" : { 
					"type" : "Point", 
					"coordinates" : [ -71.073283, 42.417500 ] 
				}
			}]
		}
	}`

	r.POST("/v1/admin/tickers/1/messages").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		SetBody(body).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
		})

	r.POST("/v1/admin/tickers/1/messages").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			type jsonResp struct {
				Data   map[string]model.Message `json:"data"`
				Status string                   `json:"status"`
				Error  interface{}              `json:"error"`
			}

			var jres jsonResp

			err := json.Unmarshal(r.Body.Bytes(), &jres)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, jres.Status)
			assert.Equal(t, nil, jres.Error)
			assert.Equal(t, 1, len(jres.Data))

			message := jres.Data["message"]

			assert.Equal(t, "message", message.Text)
			assert.Equal(t, 1, message.Ticker)
			assert.IsType(t, geojson.FeatureCollection{}, message.GeoInformation)
		})
}

func TestPostMessageWithAttachmentHandler(t *testing.T) {
	r := setup()

	ticker := model.Ticker{
		ID:       1,
		Active:   true,
		Hashtags: []string{`#hashtag`},
	}

	upload := model.Upload{
		ID:           1,
		UUID:         uuid.New().String(),
		CreationDate: time.Now(),
		TickerID:     1,
		Path:         "1/1",
		Extension:    "jpg",
		ContentType:  "image/jpeg",
	}

	storage.DB.Save(&ticker)
	storage.DB.Save(&upload)

	body := `{
		"text": "message",
		"geo_information": {
			"type" : "FeatureCollection",
			"features" : [{ 
				"type" : "Feature", 
				"properties" : {  
					"capacity" : "10", 
					"type" : "U-Rack",
					"mount" : "Surface"
				}, 
				"geometry" : { 
					"type" : "Point", 
					"coordinates" : [ -71.073283, 42.417500 ] 
				}
			}]
		},
		"attachments": [1]
	}`

	r.POST("/v1/admin/tickers/1/messages").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			type jsonResp struct {
				Data   map[string]model.MessageResponse `json:"data"`
				Status string                           `json:"status"`
				Error  interface{}                      `json:"error"`
			}

			var jres jsonResp

			err := json.Unmarshal(r.Body.Bytes(), &jres)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, jres.Status)
			assert.Equal(t, nil, jres.Error)
			assert.Equal(t, 1, len(jres.Data))

			message := jres.Data["message"]

			assert.Equal(t, "message", message.Text)
			assert.Equal(t, 1, message.Ticker)

			assert.Equal(t, 1, len(message.Attachments))
			assert.NotNil(t, message.Attachments[0].URL)
			assert.Equal(t, "image/jpeg", message.Attachments[0].ContentType)
		})
}

func TestDeleteMessageHandler(t *testing.T) {
	r := setup()

	ticker := model.Ticker{
		ID:     1,
		Active: true,
	}

	storage.DB.Save(&ticker)

	message := model.NewMessage()
	message.Text = "Text"
	message.Ticker = 1

	storage.DB.Save(message)

	r.DELETE("/v1/admin/tickers/1/messages/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
		})

	r.DELETE("/v1/admin/tickers/1/messages/2").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})

	r.DELETE("/v1/admin/tickers/1/messages/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			type jsonResp struct {
				Data   map[string]model.Message `json:"data"`
				Status string                   `json:"status"`
				Error  interface{}              `json:"error"`
			}

			var jres jsonResp

			err := json.Unmarshal(r.Body.Bytes(), &jres)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, jres.Status)
			assert.Nil(t, jres.Data)
			assert.Nil(t, jres.Error)
		})
}
