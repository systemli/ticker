package api

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

func TestGetSettingHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/settings/refresh_interval").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			assert.Equal(t, `{"data":{"setting":{"id":0,"name":"refresh_interval","value":10000}},"status":"success","error":null}`, strings.TrimSpace(r.Body.String()))
		})

	setting := model.NewSetting("refresh_interval", 20000)
	storage.DB.Save(setting)

	r.GET("/v1/admin/settings/refresh_interval").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			assert.Equal(t, `{"data":{"setting":{"id":1,"name":"refresh_interval","value":20000}},"status":"success","error":null}`, strings.TrimSpace(r.Body.String()))
		})
}

func TestGetInactiveSettingsHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/settings/inactive_settings").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			type res struct {
				Data   map[string]model.Setting `json:"data"`
				Status string                   `json:"status"`
				Error  interface{}              `json:"error"`
			}

			var response res

			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, response.Status)
			assert.Equal(t, nil, response.Error)
			assert.Equal(t, 1, len(response.Data))

			setting := response.Data["setting"]

			assert.Equal(t, 0, setting.ID)
			assert.Equal(t, "inactive_settings", setting.Name)

			value := setting.Value.(map[string]interface{})

			assert.Equal(t, model.SettingInactiveHeadline, value["headline"])
			assert.Equal(t, model.SettingInactiveSubHeadline, value["sub_headline"])
			assert.Equal(t, model.SettingInactiveDescription, value["description"])
			assert.Equal(t, model.SettingInactiveAuthor, value["author"])
			assert.Equal(t, model.SettingInactiveEmail, value["email"])
			assert.Equal(t, model.SettingInactiveHomepage, value["homepage"])
			assert.Equal(t, model.SettingInactiveTwitter, value["twitter"])
		})
}

func TestPutInactiveSettingsHandler(t *testing.T) {
	r := setup()

	body := `{
		"headline": "Headline",
		"sub_headline": "Subheadline",
		"description": "Beschreibung",
		"author": "Systemli Admin Team",
		"email": "admin@systemli.org",
		"homepage": "https://www.systemli.org",
		"twitter": "systemli"
	}`

	r.PUT("/v1/admin/settings/inactive_settings").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			type res struct {
				Data   map[string]model.Setting `json:"data"`
				Status string                   `json:"status"`
				Error  interface{}              `json:"error"`
			}

			var response res

			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, response.Status)
			assert.Equal(t, nil, response.Error)
			assert.Equal(t, 1, len(response.Data))

			setting := response.Data["setting"]

			assert.Equal(t, 1, setting.ID)
			assert.Equal(t, "inactive_settings", setting.Name)

			value := setting.Value.(map[string]interface{})

			assert.Equal(t, "Headline", value["headline"])
			assert.Equal(t, "Subheadline", value["sub_headline"])
			assert.Equal(t, "Beschreibung", value["description"])
			assert.Equal(t, "Systemli Admin Team", value["author"])
			assert.Equal(t, "admin@systemli.org", value["email"])
			assert.Equal(t, "https://www.systemli.org", value["homepage"])
			assert.Equal(t, "systemli", value["twitter"])
		})
}

func TestPutRefreshIntervalHandler(t *testing.T) {
	r := setup()

	body := `{
		"refresh_interval": 20000
	}`

	r.PUT("/v1/admin/settings/refresh_interval").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)

			type res struct {
				Data   map[string]model.Setting `json:"data"`
				Status string                   `json:"status"`
				Error  interface{}              `json:"error"`
			}

			var response res

			err := json.Unmarshal(r.Body.Bytes(), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, model.ResponseSuccess, response.Status)
			assert.Equal(t, nil, response.Error)
			assert.Equal(t, 1, len(response.Data))

			setting := response.Data["setting"]

			assert.Equal(t, 1, setting.ID)
			assert.Equal(t, "refresh_interval", setting.Name)

			value := setting.Value

			assert.Equal(t, float64(20000), value)
		})
}
