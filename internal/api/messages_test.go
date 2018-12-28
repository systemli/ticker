package api_test

import (
	"encoding/json"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/api"
	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
	"strings"
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
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			assert.Equal(t, `{"data":{"messages":[]},"status":"success","error":null}`, strings.TrimSpace(r.Body.String()))
		})

	r.GET("/v1/admin/tickers/1/messages").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1001,"message":"not found"}}`, strings.TrimSpace(r.Body.String()))
		})

	r.GET("/v1/admin/tickers/1/messages/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
		})

	message := model.NewMessage()
	message.Text = "text"
	message.Ticker = ticker.ID

	storage.DB.Save(message)

	r.GET("/v1/admin/tickers/1/messages/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
		ID:     1,
		Active: true,
		Hashtags: []string{`#hashtag`},
	}

	storage.DB.Save(&ticker)

	body := `{
		"text": "message"
	}`

	r.POST("/v1/admin/tickers/1/messages").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		SetBody(body).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
		})

	r.POST("/v1/admin/tickers/1/messages").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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

			assert.Equal(t, "message #hashtag", message.Text)
			assert.Equal(t, 1, message.Ticker)
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
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
		})

	r.DELETE("/v1/admin/tickers/1/messages/2").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})

	r.DELETE("/v1/admin/tickers/1/messages/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
