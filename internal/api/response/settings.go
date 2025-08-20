package response

import "github.com/systemli/ticker/internal/storage"

type Settings struct {
	InactiveSettings interface{} `json:"inactiveSettings,omitempty"`
}

type Setting struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func InactiveSettingsResponse(inactiveSettings storage.InactiveSettings) Setting {
	return Setting{
		Name:  storage.SettingInactiveName,
		Value: inactiveSettings,
	}
}
