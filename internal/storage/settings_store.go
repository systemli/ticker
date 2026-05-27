package storage

import (
	"encoding/json"

	"gorm.io/gorm"
)

// SettingValue is the type constraint for values stored in the settings table.
// Each is a JSON-encoded blob under a well-known name.
type SettingValue interface {
	InactiveSettings | TelegramSettings | SignalGroupSettings
}

// SettingDescriptor pairs a settings name with the function that produces the
// default value when the row is missing or unreadable.
type SettingDescriptor[T SettingValue] struct {
	Name      string
	DefaultFn func() T
}

var (
	InactiveSetting = SettingDescriptor[InactiveSettings]{
		Name:      SettingInactiveName,
		DefaultFn: DefaultInactiveSettings,
	}
	TelegramSetting = SettingDescriptor[TelegramSettings]{
		Name:      SettingTelegramName,
		DefaultFn: DefaultTelegramSettings,
	}
	SignalGroupSetting = SettingDescriptor[SignalGroupSettings]{
		Name:      SettingSignalGroupName,
		DefaultFn: DefaultSignalGroupSettings,
	}
)

// SettingsStore reads and writes raw setting rows. Use the generic
// GetSettings / SaveSettings helpers below for typed access.
type SettingsStore interface {
	GetSetting(name string) (Setting, error)
	SaveSetting(setting *Setting) error

	WithSettingsTx(tx *gorm.DB) SettingsStore
}

// GetSettings reads the row identified by desc and decodes it into T. On any
// error (missing row, bad JSON) the descriptor's default is returned.
func GetSettings[T SettingValue](s SettingsStore, desc SettingDescriptor[T]) T {
	row, err := s.GetSetting(desc.Name)
	if err != nil {
		return desc.DefaultFn()
	}
	var out T
	if err := json.Unmarshal([]byte(row.Value), &out); err != nil {
		return desc.DefaultFn()
	}
	return out
}

// SaveSettings JSON-encodes v and upserts it under desc.Name.
func SaveSettings[T SettingValue](s SettingsStore, desc SettingDescriptor[T], v T) error {
	value, err := json.Marshal(v)
	if err != nil {
		return err
	}
	row, err := s.GetSetting(desc.Name)
	if err != nil {
		row = Setting{Name: desc.Name}
	}
	row.Value = string(value)
	return s.SaveSetting(&row)
}

// GetSetting fetches a setting row by its well-known name.
func (s *SqlStorage) GetSetting(name string) (Setting, error) {
	var setting Setting
	err := s.DB.First(&setting, EqualName, name).Error
	return setting, err
}

// SaveSetting upserts a setting row.
func (s *SqlStorage) SaveSetting(setting *Setting) error {
	return s.DB.Save(setting).Error
}

// WithSettingsTx returns a SettingsStore scoped to the given transaction.
func (s *SqlStorage) WithSettingsTx(tx *gorm.DB) SettingsStore {
	return &SqlStorage{DB: tx, uploadPath: s.uploadPath}
}
