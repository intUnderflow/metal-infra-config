// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/intunderflow/metal-infra-config/entities (interfaces: Config)
//
// Generated by this command:
//
//	mockgen -destination entities/mock/mock.go github.com/intunderflow/metal-infra-config/entities Config
//

// Package mock_entities is a generated GoMock package.
package mock_entities

import (
	reflect "reflect"

	entities "github.com/intunderflow/metal-infra-config/entities"
	gomock "go.uber.org/mock/gomock"
)

// MockConfig is a mock of Config interface.
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig.
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance.
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockConfig) Delete(arg0 entities.Key) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockConfigMockRecorder) Delete(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockConfig)(nil).Delete), arg0)
}

// GetWithVersion mocks base method.
func (m *MockConfig) GetWithVersion(arg0 entities.Key) (entities.ValueWithVersion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithVersion", arg0)
	ret0, _ := ret[0].(entities.ValueWithVersion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithVersion indicates an expected call of GetWithVersion.
func (mr *MockConfigMockRecorder) GetWithVersion(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithVersion", reflect.TypeOf((*MockConfig)(nil).GetWithVersion), arg0)
}

// List mocks base method.
func (m *MockConfig) List() map[entities.Key]entities.ValueWithVersion {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].(map[entities.Key]entities.ValueWithVersion)
	return ret0
}

// List indicates an expected call of List.
func (mr *MockConfigMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockConfig)(nil).List))
}

// SetWithVersion mocks base method.
func (m *MockConfig) SetWithVersion(arg0 entities.Key, arg1 string, arg2 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWithVersion", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetWithVersion indicates an expected call of SetWithVersion.
func (mr *MockConfigMockRecorder) SetWithVersion(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWithVersion", reflect.TypeOf((*MockConfig)(nil).SetWithVersion), arg0, arg1, arg2)
}

// Sync mocks base method.
func (m *MockConfig) Sync(arg0 entities.SyncSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Sync indicates an expected call of Sync.
func (mr *MockConfigMockRecorder) Sync(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockConfig)(nil).Sync), arg0)
}
