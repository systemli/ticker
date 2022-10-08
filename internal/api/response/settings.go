package response

import "github.com/systemli/ticker/internal/storage"

type Settings struct {
	RefreshInterval  float64     `json:"refresh_interval,omitempty"`
	InactiveSettings interface{} `json:"inactive_settings,omitempty"`
}

type Setting struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func SettingResponse(setting storage.Setting) Setting {
	return Setting{
		ID:    setting.ID,
		Name:  setting.Name,
		Value: setting.Value,
	}
}
