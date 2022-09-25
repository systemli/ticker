package storage

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

type Setting struct {
	ID    int    `storm:"id,increment"`
	Name  string `storm:"unique"`
	Value interface{}
}

type InactiveSettings struct {
	Headline    string `json:"headline" binding:"required"`
	SubHeadline string `json:"sub_headline" binding:"required"`
	Description string `json:"description" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Homepage    string `json:"homepage" binding:"required"`
	Twitter     string `json:"twitter" binding:"required"`
}

type RefreshIntervalSettings struct {
	RefreshInterval int `json:"refresh_interval" binding:"required"`
}

func NewSetting(name string, value interface{}) Setting {
	return Setting{Name: name, Value: value}
}

func DefaultRefreshIntervalSetting() Setting {
	return NewSetting(SettingRefreshInterval, DefaultRefreshIntervalSettings())
}

func DefaultRefreshIntervalSettings() RefreshIntervalSettings {
	return RefreshIntervalSettings{RefreshInterval: SettingDefaultRefreshInterval}
}

func DefaultInactiveSetting() Setting {
	return NewSetting(SettingInactiveName, DefaultInactiveSettings())
}

func DefaultInactiveSettings() InactiveSettings {
	return InactiveSettings{
		Headline:    SettingInactiveHeadline,
		SubHeadline: SettingInactiveSubHeadline,
		Description: SettingInactiveDescription,
		Author:      SettingInactiveAuthor,
		Email:       SettingInactiveEmail,
		Homepage:    SettingInactiveHomepage,
		Twitter:     SettingInactiveTwitter,
	}
}
