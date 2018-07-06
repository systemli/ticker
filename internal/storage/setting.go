package storage

import (
	. "git.codecoop.org/systemli/ticker/internal/model"
)

//FindSetting lookup for a setting in storage
func FindSetting(name string) (*Setting, error) {
	var setting Setting

	err := DB.One("Name", name, &setting)
	if err != nil {
		return &setting, err
	}

	return &setting, nil
}

//GetInactiveSettings returns setting from storage or default setting
func GetInactiveSettings() *Setting {
	setting, err := FindSetting(SettingInactiveName)
	if err != nil {
		return DefaultInactiveSetting()
	}

	return setting
}

//GetRefreshInterval returns the refresh interval
func GetRefreshInterval() *Setting {
	setting, err := FindSetting(SettingRefreshInterval)
	if err != nil {
		return NewSetting(SettingRefreshInterval, SettingDefaultRefreshInterval)
	}

	return setting
}

//GetRefreshIntervalValue returns concrete integer value
func GetRefreshIntervalValue() int {
	setting := GetRefreshInterval()

	var value int
	switch sv := setting.Value.(type) {
	case float64:
		value = int(sv)
	default:
		value = sv.(int)
	}

	return value
}