package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/systemli/ticker/internal/storage"
)

func TestInactiveSettingsResponse(t *testing.T) {
	inactiveSettings := storage.DefaultInactiveSettings()

	setting := InactiveSettingsResponse(inactiveSettings)

	assert.Equal(t, storage.SettingInactiveName, setting.Name)
	assert.Equal(t, inactiveSettings, setting.Value)
}

func TestRefreshIntervalSettingsResponse(t *testing.T) {
	refreshIntervalSettings := storage.DefaultRefreshIntervalSettings()

	setting := RefreshIntervalSettingsResponse(refreshIntervalSettings)

	assert.Equal(t, storage.SettingRefreshInterval, setting.Name)
	assert.Equal(t, refreshIntervalSettings, setting.Value)
}
