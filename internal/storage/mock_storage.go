package storage

import (
	"encoding/json"
	"errors"

	"github.com/stretchr/testify/mock"
)

// MockStorage bundles the five per-aggregate mocks behind a single value for
// convenience in test setup. Field access is explicit (e.g. `m.Tickers.On(...)`)
// because each per-store mock has its own `mock.Mock`; embedding would make
// `.On` ambiguous.
type MockStorage struct {
	Users    *MockUserStore
	Tickers  *MockTickerStore
	Messages *MockMessageStore
	Uploads  *MockUploadStore
	Settings *MockSettingsStore
}

// NewMockStorage builds a MockStorage with fresh mocks for every aggregate.
func NewMockStorage() *MockStorage {
	return &MockStorage{
		Users:    &MockUserStore{},
		Tickers:  &MockTickerStore{},
		Messages: &MockMessageStore{},
		Uploads:  &MockUploadStore{},
		Settings: &MockSettingsStore{},
	}
}

// Stores returns a storage.Stores aggregate backed by the per-store mocks.
func (m *MockStorage) Stores() Stores {
	return Stores{
		Users:    m.Users,
		Tickers:  m.Tickers,
		Messages: m.Messages,
		Uploads:  m.Uploads,
		Settings: m.Settings,
	}
}

// AssertExpectations verifies all five per-store mocks at once.
func (m *MockStorage) AssertExpectations(t mock.TestingT) bool {
	ok := m.Users.AssertExpectations(t)
	ok = m.Tickers.AssertExpectations(t) && ok
	ok = m.Messages.AssertExpectations(t) && ok
	ok = m.Uploads.AssertExpectations(t) && ok
	ok = m.Settings.AssertExpectations(t) && ok
	return ok
}

// MockGetSetting wires the SettingsStore mock to return v for the given
// descriptor's name. Use this instead of poking GetSetting directly when a
// test wants the typed Get<Name>Settings() shape.
func MockGetSetting[T SettingValue](m *MockSettingsStore, desc SettingDescriptor[T], v T) *mock.Call {
	raw, _ := json.Marshal(v)
	return m.On("GetSetting", desc.Name).Return(Setting{Name: desc.Name, Value: string(raw)}, nil)
}

// MockSaveSetting wires GetSetting (missing) + SaveSetting so SaveSettings[T]
// resolves to a SaveSetting call returning the given err.
func MockSaveSetting[T SettingValue](m *MockSettingsStore, desc SettingDescriptor[T], err error) *mock.Call {
	m.On("GetSetting", desc.Name).Return(Setting{}, errors.New("not found")).Maybe()
	return m.On("SaveSetting", mock.AnythingOfType("*storage.Setting")).Return(err)
}
