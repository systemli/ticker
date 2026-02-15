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
		value := setting.Value.(TelegramSettingsValue)
		s.Equal("", value.Token)
		s.Equal("", value.BotUsername)
	})

	s.Run("with token set", func() {
		telegramSettings := storage.TelegramSettings{
			Token:       "123456789:ABCdefGHIjklMNOpqrsTUVwxyz",
			BotUsername: "test_bot",
		}

		setting := TelegramSettingsResponse(telegramSettings)

		s.Equal(storage.SettingTelegramName, setting.Name)
		value := setting.Value.(TelegramSettingsValue)
		s.Equal("****wxyz", value.Token)
		s.Equal("test_bot", value.BotUsername)
	})

	s.Run("with empty token", func() {
		telegramSettings := storage.TelegramSettings{
			Token: "",
		}

		setting := TelegramSettingsResponse(telegramSettings)

		s.Equal(storage.SettingTelegramName, setting.Name)
		value := setting.Value.(TelegramSettingsValue)
		s.Equal("", value.Token)
		s.Equal("", value.BotUsername)
	})
}

func (s *SettingsResponseTestSuite) TestMaskToken() {
	s.Run("empty token", func() {
		s.Equal("", maskToken(""))
	})

	s.Run("short token", func() {
		s.Equal("****", maskToken("ab"))
	})

	s.Run("normal token", func() {
		s.Equal("****wxyz", maskToken("123456789:ABCdefGHIjklMNOpqrsTUVwxyz"))
	})
}

func TestSettingsResponseTestSuite(t *testing.T) {
	suite.Run(t, new(SettingsResponseTestSuite))
}
