package response

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/storage"
)

type SettingsResponseTestSuite struct {
	suite.Suite
}

func (s *SettingsResponseTestSuite) TestInactiveSettingsResponse() {
	inactiveSettings := storage.DefaultInactiveSettings()

	setting := InactiveSettingsResponse(inactiveSettings)

	s.Equal(storage.SettingInactiveName, setting.Name)
	s.Equal(inactiveSettings, setting.Value)
}

func (s *SettingsResponseTestSuite) TestTelegramSettingsResponse() {
	s.Run("with default settings", func() {
		telegramSettings := storage.DefaultTelegramSettings()

		setting := TelegramSettingsResponse(telegramSettings)

		s.Equal(storage.SettingTelegramName, setting.Name)
		s.Equal(telegramSettings, setting.Value)
		s.Equal("", setting.Value.(storage.TelegramSettings).Token)
	})

	s.Run("with token set", func() {
		telegramSettings := storage.TelegramSettings{
			Token: "123456789:ABCdefGHIjklMNOpqrsTUVwxyz",
		}

		setting := TelegramSettingsResponse(telegramSettings)

		s.Equal(storage.SettingTelegramName, setting.Name)
		s.Equal(telegramSettings, setting.Value)
		s.Equal("123456789:ABCdefGHIjklMNOpqrsTUVwxyz", setting.Value.(storage.TelegramSettings).Token)
	})

	s.Run("with empty token", func() {
		telegramSettings := storage.TelegramSettings{
			Token: "",
		}

		setting := TelegramSettingsResponse(telegramSettings)

		s.Equal(storage.SettingTelegramName, setting.Name)
		s.Equal(telegramSettings, setting.Value)
		s.Equal("", setting.Value.(storage.TelegramSettings).Token)
	})

	s.Run("response structure validation", func() {
		telegramSettings := storage.TelegramSettings{
			Token: "test-token",
		}

		setting := TelegramSettingsResponse(telegramSettings)

		s.IsType(Setting{}, setting)
		s.IsType(storage.TelegramSettings{}, setting.Value)
		s.NotEmpty(setting.Name)
	})
}

func TestSettingsResponseTestSuite(t *testing.T) {
	suite.Run(t, new(SettingsResponseTestSuite))
}
