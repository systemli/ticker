package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettingConstants(t *testing.T) {
	t.Run("setting names are defined correctly", func(t *testing.T) {
		assert.Equal(t, "inactive_settings", SettingInactiveName)
		assert.Equal(t, "telegram_settings", SettingTelegramName)
		assert.Equal(t, "signal_group_settings", SettingSignalGroupName)
		assert.NotEqual(t, SettingInactiveName, SettingTelegramName)
		assert.NotEqual(t, SettingTelegramName, SettingSignalGroupName)
		assert.NotEqual(t, SettingInactiveName, SettingSignalGroupName)
	})
}

func TestDefaultInactiveSettings(t *testing.T) {
	t.Run("returns correct default values", func(t *testing.T) {
		defaults := DefaultInactiveSettings()

		assert.Equal(t, SettingInactiveHeadline, defaults.Headline)
		assert.Equal(t, SettingInactiveSubHeadline, defaults.SubHeadline)
		assert.Equal(t, SettingInactiveDescription, defaults.Description)
		assert.Equal(t, SettingInactiveAuthor, defaults.Author)
		assert.Equal(t, SettingInactiveEmail, defaults.Email)
		assert.Equal(t, SettingInactiveHomepage, defaults.Homepage)
		assert.Equal(t, SettingInactiveTwitter, defaults.Twitter)
	})
}

func TestDefaultTelegramSettings(t *testing.T) {
	t.Run("returns correct default values", func(t *testing.T) {
		defaults := DefaultTelegramSettings()

		assert.Equal(t, "", defaults.Token)
		assert.Equal(t, "", defaults.BotUsername)
	})
}

func TestDefaultSignalGroupSettings(t *testing.T) {
	t.Run("returns correct default values", func(t *testing.T) {
		defaults := DefaultSignalGroupSettings()

		assert.Equal(t, "", defaults.ApiUrl)
		assert.Equal(t, "", defaults.Account)
		assert.Equal(t, "", defaults.Avatar)
		assert.False(t, defaults.Enabled())
	})
}

func TestSignalGroupSettingsEnabled(t *testing.T) {
	t.Run("returns false when both fields are empty", func(t *testing.T) {
		settings := SignalGroupSettings{}
		assert.False(t, settings.Enabled())
	})

	t.Run("returns false when only ApiUrl is set", func(t *testing.T) {
		settings := SignalGroupSettings{ApiUrl: "http://localhost:8080"}
		assert.False(t, settings.Enabled())
	})

	t.Run("returns false when only Account is set", func(t *testing.T) {
		settings := SignalGroupSettings{Account: "+491234567890"}
		assert.False(t, settings.Enabled())
	})

	t.Run("returns true when both ApiUrl and Account are set", func(t *testing.T) {
		settings := SignalGroupSettings{ApiUrl: "http://localhost:8080", Account: "+491234567890"}
		assert.True(t, settings.Enabled())
	})
}
