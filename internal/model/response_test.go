package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
)

func TestNewJSONSuccessResponse(t *testing.T) {
	d := []string{"value1", "value2"}
	r := model.NewJSONSuccessResponse("user", d)

	assert.Equal(t, "success", r.Status)
	assert.Equal(t, map[string]interface{}{"user": []string{"value1", "value2"}}, r.Data)
	assert.Nil(t, r.Error)
}

func TestNewJSONErrorResponse(t *testing.T) {
	r := model.NewJSONErrorResponse(1, "error")

	assert.Equal(t, "error", r.Status)
	assert.Equal(t, map[string]interface{}{"code": 1, "message": "error"}, r.Error)
	assert.Equal(t, map[string]interface{}{}, r.Data)
}
