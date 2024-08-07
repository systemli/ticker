// Code generated by mockery v2.43.0. DO NOT EDIT.

package bridge

import (
	mock "github.com/stretchr/testify/mock"
	storage "github.com/systemli/ticker/internal/storage"
)

// MockBridge is an autogenerated mock type for the Bridge type
type MockBridge struct {
	mock.Mock
}

func (_m *MockBridge) Update(ticker storage.Ticker) error {
	ret := _m.Called(ticker)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(storage.Ticker) error); ok {
		r0 = rf(ticker)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ticker, message
func (_m *MockBridge) Delete(ticker storage.Ticker, message *storage.Message) error {
	ret := _m.Called(ticker, message)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(storage.Ticker, *storage.Message) error); ok {
		r0 = rf(ticker, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Send provides a mock function with given fields: ticker, message
func (_m *MockBridge) Send(ticker storage.Ticker, message *storage.Message) error {
	ret := _m.Called(ticker, message)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(storage.Ticker, *storage.Message) error); ok {
		r0 = rf(ticker, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockBridge creates a new instance of MockBridge. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockBridge(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockBridge {
	mock := &MockBridge{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
