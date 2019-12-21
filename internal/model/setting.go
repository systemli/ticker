package model

const (
	SettingInactiveName           = `inactive_settings`
	SettingRefreshInterval        = `refresh_interval`
	SettingInactiveHeadline       = `The ticker is currently inactive.`
	SettingInactiveSubHeadline    = `Please contact us if you want to use it.`
	SettingInactiveDescription    = `...`
	SettingInactiveAuthor         = `systemli.org Ticker Team`
	SettingInactiveEmail          = `admin@systemli.org`
	SettingInactiveHomepage       = `https://www.systemli.org/`
	SettingInactiveTwitter        = `systemli`
	SettingDefaultRefreshInterval = 10000
)

//Setting represents an global setting for ticker
type Setting struct {
	ID    int    `storm:"id,increment"`
	Name  string `storm:"unique"`
	Value interface{}
}

//InactiveSetting represents the inactive properties
type InactiveSettings struct {
	Headline    string `json:"headline,omitempty"`
	SubHeadline string `json:"sub_headline,omitempty"`
	Description string `json:"description,omitempty"`
	Author      string `json:"author,omitempty"`
	Email       string `json:"email,omitempty"`
	Homepage    string `json:"homepage,omitempty"`
	Twitter     string `json:"twitter,omitempty"`
}

type SettingResponse struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

//NewSetting return a new setting instance
func NewSetting(name string, value interface{}) *Setting {
	return &Setting{Name: name, Value: value}
}

//DefaultInactiveSetting returns default inactive_settings
func DefaultInactiveSetting() *Setting {
	return NewSetting(SettingInactiveName, DefaultInactiveSettings())
}

//DefaultInactiveSettings returns the default value for inactive_settings
func DefaultInactiveSettings() *InactiveSettings {
	return &InactiveSettings{
		Headline:    SettingInactiveHeadline,
		SubHeadline: SettingInactiveSubHeadline,
		Description: SettingInactiveDescription,
		Author:      SettingInactiveAuthor,
		Email:       SettingInactiveEmail,
		Homepage:    SettingInactiveHomepage,
		Twitter:     SettingInactiveTwitter,
	}
}

//NewSettingResponse returns a setting response
func NewSettingResponse(setting *Setting) *SettingResponse {
	return &SettingResponse{
		ID:    setting.ID,
		Name:  setting.Name,
		Value: setting.Value,
	}
}
