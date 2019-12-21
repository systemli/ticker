package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
)

func TestNewSetting(t *testing.T) {
	s := model.NewSetting("key", "value")

	assert.Equal(t, "key", s.Name)
	assert.Equal(t, "value", s.Value)
}

func TestDefaultInactiveSetting(t *testing.T) {
	d := model.DefaultInactiveSetting()

	assert.Equal(t, "inactive_settings", d.Name)
	assert.NotNil(t, d.Value)
}

func TestNewSettingResponse(t *testing.T) {
	s := model.NewSetting("key", "value")
	r := model.NewSettingResponse(s)

	assert.Equal(t, "key", r.Name)
	assert.Equal(t, "value", r.Value)
}
