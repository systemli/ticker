package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettingConstants(t *testing.T) {
	t.Run("setting names are defined correctly", func(t *testing.T) {
		assert.Equal(t, "inactive_settings", SettingInactiveName)
		assert.Equal(t, "telegram_settings", SettingTelegramName)
		assert.NotEqual(t, SettingInactiveName, SettingTelegramName)
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
