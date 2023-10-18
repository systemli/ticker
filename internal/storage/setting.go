package storage

import "time"

const (
	SettingInactiveName               = `inactive_settings`
	SettingRefreshInterval            = `refresh_interval`
	SettingInactiveHeadline           = `The ticker is currently inactive.`
	SettingInactiveSubHeadline        = `Please contact us if you want to use it.`
	SettingInactiveDescription        = `...`
	SettingInactiveAuthor             = `systemli.org Ticker Team`
	SettingInactiveEmail              = `admin@systemli.org`
	SettingInactiveHomepage           = `https://www.systemli.org/`
	SettingInactiveTwitter            = `systemli`
	SettingDefaultRefreshInterval int = 10000
)

type Setting struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"unique"`
	Value     string `gorm:"type:json"`
}

type InactiveSettings struct {
	Headline    string `json:"headline" binding:"required"`
	SubHeadline string `json:"subHeadline" binding:"required"`
	Description string `json:"description" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Homepage    string `json:"homepage" binding:"required"`
	Twitter     string `json:"twitter" binding:"required"`
}

type RefreshIntervalSettings struct {
	RefreshInterval int `json:"refreshInterval" binding:"required"`
}

func DefaultRefreshIntervalSettings() RefreshIntervalSettings {
	return RefreshIntervalSettings{RefreshInterval: SettingDefaultRefreshInterval}
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
