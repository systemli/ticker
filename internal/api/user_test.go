package api_test

import (
	"testing"
	"github.com/appleboy/gofight"
	"strings"
	"git.codecoop.org/systemli/ticker/internal/api"
	"github.com/stretchr/testify/assert"
	"git.codecoop.org/systemli/ticker/internal/model"
	"encoding/json"
	"git.codecoop.org/systemli/ticker/internal/storage"
	"fmt"
)

func TestGetUsers(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/users").
		SetHeader(map[string]string{"Authorization": "Bearer " + Token}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)
		assert.Equal(t, `{"data":{"users":null},"status":"success","error":null}`, strings.TrimSpace(r.Body.String()))
	})
}

func TestGetUser(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/users/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + Token}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1001,"message":"not found"}}`, strings.TrimSpace(r.Body.String()))
	})
}

func TestPostUser(t *testing.T) {
	r := setup()

	body := `{
		"email": "louis@systemli.org",
		"password": "password12"
	}`

	r.POST("/v1/admin/users").
		SetHeader(map[string]string{"Authorization": "Bearer " + Token}).
		SetBody(body).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)

		var response struct {
			Data   map[string]model.User `json:"data"`
			Status string                `json:"status"`
			Error  interface{}           `json:"error"`
		}

		err := json.Unmarshal(r.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.ResponseSuccess, response.Status)
		assert.Equal(t, nil, response.Error)
		assert.Equal(t, 1, len(response.Data))
		assert.Equal(t, "louis@systemli.org", response.Data["user"].Email)
	})
}

func TestPutUser(t *testing.T) {
	r := setup()

	u, err := model.NewUser("louis@systemli.org", "password")
	if err != nil {
		t.Fail()
	}

	storage.DB.Save(u)

	body := `{
		"email": "admin@systemli.org",
		"password": "password13",
		"role": "user",
		"is_super_admin": true,
		"tickers": [1,2,3]
	}`

	r.PUT(fmt.Sprintf(`/v1/admin/users/%d`, u.ID)).
		SetHeader(map[string]string{"Authorization": "Bearer " + Token}).
		SetBody(body).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)

		var response struct {
			Data   map[string]model.User `json:"data"`
			Status string                `json:"status"`
			Error  interface{}           `json:"error"`
		}

		err := json.Unmarshal(r.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.ResponseSuccess, response.Status)
		assert.Equal(t, nil, response.Error)
		assert.Equal(t, 1, len(response.Data))
		assert.Equal(t, u.ID, response.Data["user"].ID)
		assert.Equal(t, "admin@systemli.org", response.Data["user"].Email)
		assert.Equal(t, "user", response.Data["user"].Role)

		var user model.User
		err = storage.DB.One("ID", u.ID, &user)
		if err != nil {
			t.Fail()
		}

		assert.NotEmpty(t, user.EncryptedPassword)
		assert.Equal(t, true, user.IsSuperAdmin)
		assert.Equal(t, []int{1, 2, 3}, user.Tickers)
	})
}

func TestDeleteUser(t *testing.T) {
	r := setup()

	user := model.User{
		ID:     1,
	}

	storage.DB.Save(&user)

	r.DELETE("/v1/admin/users/2").
		SetHeader(map[string]string{"Authorization": "Bearer " + Token}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
	})

	r.DELETE("/v1/admin/users/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + Token}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
}