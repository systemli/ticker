package storage

import (
	"github.com/stretchr/testify/mock"
)

// Convenience helpers attached to MockSettingsStore so tests can keep the
// previous Get<Name>Settings / Save<Name>Settings call shape after the split
// of the storage interface into per-aggregate stores.

func (m *MockSettingsStore) MockGetInactive(v InactiveSettings) *mock.Call {
	return MockGetSetting(m, InactiveSetting, v)
}

func (m *MockSettingsStore) MockGetTelegram(v TelegramSettings) *mock.Call {
	return MockGetSetting(m, TelegramSetting, v)
}

func (m *MockSettingsStore) MockGetSignalGroup(v SignalGroupSettings) *mock.Call {
	return MockGetSetting(m, SignalGroupSetting, v)
}

func (m *MockSettingsStore) MockSaveInactive(err error) *mock.Call {
	return MockSaveSetting(m, InactiveSetting, err)
}

func (m *MockSettingsStore) MockSaveTelegram(err error) *mock.Call {
	return MockSaveSetting(m, TelegramSetting, err)
}

func (m *MockSettingsStore) MockSaveSignalGroup(err error) *mock.Call {
	return MockSaveSetting(m, SignalGroupSetting, err)
}
