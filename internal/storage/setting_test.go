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

	t.Run("returns consistent values on multiple calls", func(t *testing.T) {
		defaults1 := DefaultInactiveSettings()
		defaults2 := DefaultInactiveSettings()

		assert.Equal(t, defaults1, defaults2)
	})
}

func TestDefaultTelegramSettings(t *testing.T) {
	t.Run("returns correct default values", func(t *testing.T) {
		defaults := DefaultTelegramSettings()

		assert.Equal(t, "", defaults.Token)
	})

	t.Run("returns consistent values on multiple calls", func(t *testing.T) {
		defaults1 := DefaultTelegramSettings()
		defaults2 := DefaultTelegramSettings()

		assert.Equal(t, defaults1, defaults2)
	})

	t.Run("default token is empty string", func(t *testing.T) {
		defaults := DefaultTelegramSettings()

		assert.Empty(t, defaults.Token)
		assert.NotNil(t, defaults.Token) // Should be empty string, not nil
	})
}

func TestTelegramSettingsStruct(t *testing.T) {
	t.Run("can be created with token", func(t *testing.T) {
		settings := TelegramSettings{
			Token: "123456789:ABCdefGHIjklMNOpqrsTUVwxyz",
		}

		assert.Equal(t, "123456789:ABCdefGHIjklMNOpqrsTUVwxyz", settings.Token)
	})

	t.Run("can be created with empty token", func(t *testing.T) {
		settings := TelegramSettings{
			Token: "",
		}

		assert.Equal(t, "", settings.Token)
	})

	t.Run("zero value has empty token", func(t *testing.T) {
		var settings TelegramSettings

		assert.Equal(t, "", settings.Token)
	})
}

func TestInactiveSettingsStruct(t *testing.T) {
	t.Run("can be created with all fields", func(t *testing.T) {
		settings := InactiveSettings{
			Headline:    "Test Headline",
			SubHeadline: "Test SubHeadline",
			Description: "Test Description",
			Author:      "Test Author",
			Email:       "test@example.com",
			Homepage:    "https://test.example",
			Twitter:     "test_twitter",
		}

		assert.Equal(t, "Test Headline", settings.Headline)
		assert.Equal(t, "Test SubHeadline", settings.SubHeadline)
		assert.Equal(t, "Test Description", settings.Description)
		assert.Equal(t, "Test Author", settings.Author)
		assert.Equal(t, "test@example.com", settings.Email)
		assert.Equal(t, "https://test.example", settings.Homepage)
		assert.Equal(t, "test_twitter", settings.Twitter)
	})

	t.Run("zero value has empty fields", func(t *testing.T) {
		var settings InactiveSettings

		assert.Equal(t, "", settings.Headline)
		assert.Equal(t, "", settings.SubHeadline)
		assert.Equal(t, "", settings.Description)
		assert.Equal(t, "", settings.Author)
		assert.Equal(t, "", settings.Email)
		assert.Equal(t, "", settings.Homepage)
		assert.Equal(t, "", settings.Twitter)
	})
}
