package response

import "github.com/systemli/ticker/internal/storage"

type Settings struct {
	InactiveSettings interface{} `json:"inactiveSettings,omitempty"`
}

type Setting struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type TelegramSettingsValue struct {
	Token       string `json:"token"`
	BotUsername string `json:"botUsername"`
}

func InactiveSettingsResponse(inactiveSettings storage.InactiveSettings) Setting {
	return Setting{
		Name:  storage.SettingInactiveName,
		Value: inactiveSettings,
	}
}

func TelegramSettingsResponse(telegramSettings storage.TelegramSettings) Setting {
	return Setting{
		Name: storage.SettingTelegramName,
		Value: TelegramSettingsValue{
			Token:       maskToken(telegramSettings.Token),
			BotUsername: telegramSettings.BotUsername,
		},
	}
}

// maskToken returns a masked version of the token, showing only the last 4 characters.
// If the token is empty, it returns an empty string.
func maskToken(token string) string {
	if token == "" {
		return ""
	}
	if len(token) <= 4 {
		return "****"
	}
	return "****" + token[len(token)-4:]
}
