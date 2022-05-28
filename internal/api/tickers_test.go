package api

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

var AdminToken string
var UserToken string

func TestGetTickersHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/tickers").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			assert.Equal(t, `{"data":{"tickers":null},"status":"success","error":null}`, strings.TrimSpace(r.Body.String()))
		})

	r.GET("/v1/admin/tickers").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			assert.Equal(t, `{"data":{"tickers":null},"status":"success","error":null}`, strings.TrimSpace(r.Body.String()))
		})
}

func TestGetTickerHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1001,"message":"not found"}}`, strings.TrimSpace(r.Body.String()))
		})

	r.GET("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
		})
}

func TestPostTickerHandler(t *testing.T) {
	r := setup()

	body := `{
		"title": "Ticker",
		"domain": "prozessticker.org",
		"description": "Beschreibung",
		"active": true,
		"hashtags": ["#test"],
		"information": {
			"url": "https://www.systemli.org",
			"email": "admin@systemli.org",
			"twitter": "systemli",
			"telegram": "https://t.me/bla"
		},
		"location": {
			"lat": 1.1,
			"lon": 2.2
		}
	}`

	r.POST("/v1/admin/tickers").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			type jsonResp struct {
				Data   map[string]model.Ticker `json:"data"`
				Status string                  `json:"status"`
				Error  interface{}             `json:"error"`
			}

			var jres jsonResp

			err := json.Unmarshal(r.Body.Bytes(), &jres)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, jres.Status)
			assert.Equal(t, nil, jres.Error)
			assert.Equal(t, 1, len(jres.Data))

			ticker := jres.Data["ticker"]

			assert.Equal(t, "Ticker", ticker.Title)
			assert.Equal(t, "prozessticker.org", ticker.Domain)
			assert.Equal(t, true, ticker.Active)
			assert.Equal(t, []string{"#test"}, ticker.Hashtags)
			assert.Equal(t, "https://www.systemli.org", ticker.Information.URL)
			assert.Equal(t, "admin@systemli.org", ticker.Information.Email)
			assert.Equal(t, "systemli", ticker.Information.Twitter)
			assert.Equal(t, "https://t.me/bla", ticker.Information.Telegram)
			assert.Equal(t, 1.1, ticker.Location.Lat)
			assert.Equal(t, 2.2, ticker.Location.Lon)
		})

	r.POST("/v1/admin/tickers").
		SetBody(body).
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
		})
}

func TestPutTickerHandler(t *testing.T) {
	r := setup()

	ticker := model.Ticker{
		ID:       1,
		Active:   true,
		Hashtags: []string{"test"},
		Domain:   "demoticker.org",
	}

	storage.DB.Save(&ticker)

	body := `{
		"title": "Ticker",
		"domain": "prozessticker.org",
		"description": "Beschreibung",
		"active": false,
		"hashtags": [],
		"information": {
			"url": "https://www.systemli.org",
			"email": "admin@systemli.org",
			"telegram": "https://t.me/bla"
		},
		"location": {
			"lat": 1.1,
			"lon": 2.2
		}
	}`

	r.PUT("/v1/admin/tickers/100").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})

	r.PUT("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(`malicious data`).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})

	r.PUT("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			type jsonResp struct {
				Data   map[string]model.Ticker `json:"data"`
				Status string                  `json:"status"`
				Error  interface{}             `json:"error"`
			}

			var jres jsonResp

			err := json.Unmarshal(r.Body.Bytes(), &jres)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, jres.Status)
			assert.Equal(t, nil, jres.Error)
			assert.Equal(t, 1, len(jres.Data))

			ticker := jres.Data["ticker"]

			assert.Equal(t, 1, ticker.ID)
			assert.Equal(t, "Ticker", ticker.Title)
			assert.Equal(t, "prozessticker.org", ticker.Domain)
			assert.Equal(t, false, ticker.Active)
			assert.Equal(t, []string{}, ticker.Hashtags)
			assert.Equal(t, 1.1, ticker.Location.Lat)
			assert.Equal(t, 2.2, ticker.Location.Lon)
			assert.Equal(t, "https://t.me/bla", ticker.Information.Telegram)
		})

	r.PUT("/v1/admin/tickers/1").
		SetBody(body).
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
		})

	body = `{
		"title": "Ticker",
		"domain": "prozessticker.org",
		"description": "Beschreibung",
		"active": false,
		"hashtags": [],
		"information": {
			"url": "https://www.systemli.org",
			"email": "admin@systemli.org"
		},
		"location": {
			"lat": 0,
			"lon": 0
		}
	}`

	r.PUT("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			type jsonResp struct {
				Data   map[string]model.Ticker `json:"data"`
				Status string                  `json:"status"`
				Error  interface{}             `json:"error"`
			}

			var jres jsonResp

			err := json.Unmarshal(r.Body.Bytes(), &jres)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, jres.Status)
			assert.Equal(t, nil, jres.Error)
			assert.Equal(t, 1, len(jres.Data))

			ticker := jres.Data["ticker"]

			assert.Equal(t, 1, ticker.ID)
			assert.Equal(t, "Ticker", ticker.Title)
			assert.Equal(t, "prozessticker.org", ticker.Domain)
			assert.Equal(t, false, ticker.Active)
			assert.Equal(t, []string{}, ticker.Hashtags)
			assert.Equal(t, 0.0, ticker.Location.Lat)
			assert.Equal(t, 0.0, ticker.Location.Lon)
		})
}

func TestPutTickerTelegramHandler(t *testing.T) {
	r := setup()

	ticker := model.Ticker{
		ID:     1,
		Active: true,
		Domain: "demoticker.org",
	}

	storage.DB.Save(&ticker)

	r.PUT("/v1/admin/tickers/1/telegram").
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 401, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1002,"message":"auth header is empty"}}`, strings.TrimSpace(r.Body.String()))
		})

	r.PUT("/v1/admin/tickers/a/telegram").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})

	r.PUT("/v1/admin/tickers/2/telegram").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})

	r.PUT("/v1/admin/tickers/1/telegram").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
		})

	body := `{
		"active": true,
		"token": "",
		"channel_name": "@channel_name"
	}`

	r.PUT("/v1/admin/tickers/1/telegram").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
		})
}

func TestDeleteTickerHandler(t *testing.T) {
	r := setup()

	ticker := model.Ticker{
		ID:     1,
		Active: true,
	}

	storage.DB.Save(&ticker)

	r.DELETE("/v1/admin/tickers/2").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})

	r.DELETE("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			var jres struct {
				Data   map[string]model.Message `json:"data"`
				Status string                   `json:"status"`
				Error  interface{}              `json:"error"`
			}

			err := json.Unmarshal(r.Body.Bytes(), &jres)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, jres.Status)
			assert.Nil(t, jres.Data)
			assert.Nil(t, jres.Error)
		})

	r.DELETE("/v1/admin/tickers/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
		})
}

func TestResetTickerHandler(t *testing.T) {
	r := setup()

	ticker := model.Ticker{
		ID:     1,
		Active: true,
		Twitter: model.Twitter{
			Token:  "token",
			Secret: "secret",
			Active: true,
		},
		Location: model.Location{
			Lat: 1.1,
			Lon: 2.2,
		},
	}

	storage.DB.Save(&ticker)

	message := model.NewMessage()
	message.Text = "Text"
	message.Ticker = 1

	storage.DB.Save(message)

	r.PUT("/v1/admin/tickers/2/reset").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})

	r.PUT("/v1/admin/tickers/1/reset").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
			assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
		})

	r.PUT("/v1/admin/tickers/1/reset").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			type jsonResp struct {
				Data   map[string]model.Ticker `json:"data"`
				Status string                  `json:"status"`
				Error  interface{}             `json:"error"`
			}

			var jres jsonResp

			err := json.Unmarshal(r.Body.Bytes(), &jres)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, jres.Status)
			assert.Equal(t, nil, jres.Error)
			assert.Equal(t, 1, len(jres.Data))

			ticker := jres.Data["ticker"]

			assert.Equal(t, 1, ticker.ID)
			assert.Equal(t, false, ticker.Active)
			assert.Equal(t, model.Location{}, ticker.Location)

			cnt, err := storage.DB.Count(model.NewMessage())
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, 0, cnt)
		})
}

func TestGetTickerUsersHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/tickers/2/users").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})

	ticker := model.Ticker{
		ID:     1,
		Active: true,
	}

	storage.DB.Save(&ticker)

	r.GET("/v1/admin/tickers/1/users").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
		})

	r.GET("/v1/admin/tickers/1/users").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
		})

	user, _ := model.NewUser("user@systemli.org", "password")
	user.Tickers = []int{ticker.ID}

	storage.DB.Save(user)

	r.GET("/v1/admin/tickers/1/users").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			type jsonResp struct {
				Data   map[string][]model.UserResponse `json:"data"`
				Status string                          `json:"status"`
				Error  interface{}                     `json:"error"`
			}

			var response jsonResp

			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, response.Status)
			assert.Equal(t, nil, response.Error)
			assert.Equal(t, 1, len(response.Data))

			assert.Equal(t, 1, len(response.Data["users"]))
			assert.Equal(t, "user@systemli.org", response.Data["users"][0].Email)
		})
}

func TestPutTickerUsersHandler(t *testing.T) {
	r := setup()

	r.PUT("/v1/admin/tickers/2/users").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})

	ticker := model.Ticker{
		ID:     1,
		Active: true,
	}

	storage.DB.Save(&ticker)

	r.PUT("/v1/admin/tickers/1/users").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
		})

	body := `{"users": [2]}`

	r.PUT("/v1/admin/tickers/1/users").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			type jsonResp struct {
				Data   map[string][]model.UserResponse `json:"data"`
				Status string                          `json:"status"`
				Error  interface{}                     `json:"error"`
			}

			var response jsonResp

			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, response.Status)
			assert.Equal(t, nil, response.Error)
			assert.Equal(t, 1, len(response.Data))

			assert.Equal(t, 1, len(response.Data["users"]))
			assert.Equal(t, "louis@systemli.org", response.Data["users"][0].Email)
		})
}

func TestDeleteTickerUserHandler(t *testing.T) {
	r := setup()

	r.DELETE("/v1/admin/tickers/1/users/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})

	ticker := model.Ticker{
		ID:     1,
		Active: true,
	}

	storage.DB.Save(&ticker)

	r.DELETE("/v1/admin/tickers/1/users/10").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})

	user := model.User{
		ID:      10,
		Email:   "user_10@systemli.org",
		Tickers: []int{1},
	}

	storage.DB.Save(&user)

	r.DELETE("/v1/admin/tickers/1/users/10").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			type jsonResp struct {
				Data   map[string][]model.UserResponse `json:"data"`
				Status string                          `json:"status"`
				Error  interface{}                     `json:"error"`
			}

			var response jsonResp

			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, response.Status)
			assert.Equal(t, nil, response.Error)
			assert.Equal(t, 1, len(response.Data))
			assert.Equal(t, 0, len(response.Data["users"]))
		})

}

func setup() *gofight.RequestConfig {
	gin.SetMode(gin.TestMode)

	model.Config = model.NewConfig()
	model.Config.UploadPath = os.TempDir()

	if storage.DB == nil {
		storage.DB = storage.OpenDB(fmt.Sprintf("%s/ticker_%d.db", os.TempDir(), time.Now().Unix()))
	}
	storage.DB.Drop("Ticker")
	storage.DB.Drop("Message")
	storage.DB.Drop("User")
	storage.DB.Drop("Setting")

	admin, _ := model.NewUser("admin@systemli.org", "password")
	admin.IsSuperAdmin = true

	storage.DB.Save(admin)

	user, _ := model.NewUser("louis@systemli.org", "password")
	storage.DB.Save(user)

	if AdminToken == "" {
		AdminToken = token("admin@systemli.org", "password")
	}
	if UserToken == "" {
		UserToken = token("louis@systemli.org", "password")
	}

	return gofight.New()
}

func token(username, password string) string {
	var token string

	r := gofight.New()
	r.POST("/v1/admin/login").
		SetBody(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password)).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			var response struct {
				Code   int       `json:"code"`
				Expire time.Time `json:"expire"`
				Token  string    `json:"token"`
			}

			json.Unmarshal(r.Body.Bytes(), &response)

			token = response.Token
		})

	return token
}
