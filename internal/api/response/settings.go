package response

import "github.com/systemli/ticker/internal/storage"

type Settings struct {
	RefreshInterval  int         `json:"refresh_interval,omitempty"`
	InactiveSettings interface{} `json:"inactive_settings,omitempty"`
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

func RefreshIntervalSettingsResponse(refreshIntervalSettings storage.RefreshIntervalSettings) Setting {
	return Setting{
		Name:  storage.SettingRefreshInterval,
		Value: refreshIntervalSettings,
	}
}
