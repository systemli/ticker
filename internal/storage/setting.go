package storage

import . "git.codecoop.org/systemli/ticker/internal/model"

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