package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

func TestFindSetting(t *testing.T) {
	setup()

	_, err := storage.FindSetting("key")
	if err == nil {
		t.Fail()
		return
	}

	setting := model.NewSetting("key", "value")
	_ = storage.DB.Save(setting)

	found, err := storage.FindSetting("key")
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 1, found.ID)
	assert.Equal(t, "value", found.Value)
}

func TestGetInactiveSettings(t *testing.T) {
	setup()

	s := storage.GetInactiveSettings()

	assert.Equal(t, 0, s.ID)

	_ = storage.DB.Save(s)

	s = storage.GetInactiveSettings()

	assert.Equal(t, 1, s.ID)
}

func TestGetRefreshInterval(t *testing.T) {
	setup()

	s := storage.GetRefreshInterval()

	assert.Equal(t, 0, s.ID)

	_ = storage.DB.Save(s)

	s = storage.GetRefreshInterval()

	assert.Equal(t, 1, s.ID)
}

func TestGetRefreshIntervalValue(t *testing.T) {
	setup()

	v := storage.GetRefreshIntervalValue()

	assert.Equal(t, 10000, v)

	var value float64
	value = 20000.00
	s := model.NewSetting(model.SettingRefreshInterval, value)
	_ = storage.DB.Save(s)

	v = storage.GetRefreshIntervalValue()

	assert.Equal(t, 20000, v)
}
