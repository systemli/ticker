package api_test

import (
	"strings"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"

	"git.codecoop.org/systemli/ticker/internal/api"
	"git.codecoop.org/systemli/ticker/internal/model"
	"git.codecoop.org/systemli/ticker/internal/storage"
	"encoding/json"
)

func TestGetSettingHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/settings/refresh_interval").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1001,"message":"setting not found"}}`, strings.TrimSpace(r.Body.String()))
	})

	setting := model.NewSetting("refresh_interval", 10000)
	storage.DB.Save(setting)

	r.GET("/v1/admin/settings/refresh_interval").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)
		assert.Equal(t, `{"data":{"setting":{"id":1,"name":"refresh_interval","value":10000}},"status":"success","error":null}`, strings.TrimSpace(r.Body.String()))
	})
}

func TestGetInactiveSettingsHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/settings/inactive_settings").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
