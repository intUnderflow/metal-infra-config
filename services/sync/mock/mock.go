// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/intunderflow/metal-infra-config/services/sync (interfaces: Sync,RPC)
//
// Generated by this command:
//
//	mockgen -destination services/sync/mock/mock.go github.com/intunderflow/metal-infra-config/services/sync Sync,RPC
//

// Package mock_sync is a generated GoMock package.
package mock_sync

import (
	context "context"
	reflect "reflect"
	time "time"

	entities "github.com/intunderflow/metal-infra-config/entities"
	closeable "github.com/intunderflow/metal-infra-config/pkg/closeable"
	proto "github.com/intunderflow/metal-infra-config/proto"
	gomock "go.uber.org/mock/gomock"
)

// MockSync is a mock of Sync interface.
type MockSync struct {
	ctrl     *gomock.Controller
	recorder *MockSyncMockRecorder
}

// MockSyncMockRecorder is the mock recorder for MockSync.
type MockSyncMockRecorder struct {
	mock *MockSync
}

// NewMockSync creates a new mock instance.
func NewMockSync(ctrl *gomock.Controller) *MockSync {
	mock := &MockSync{ctrl: ctrl}
	mock.recorder = &MockSyncMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSync) EXPECT() *MockSyncMockRecorder {
	return m.recorder
}

// Sync mocks base method.
func (m *MockSync) Sync(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Sync indicates an expected call of Sync.
func (mr *MockSyncMockRecorder) Sync(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockSync)(nil).Sync), arg0)
}

// SyncPeriodically mocks base method.
func (m *MockSync) SyncPeriodically(arg0 context.Context, arg1 time.Duration) closeable.Closeable {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncPeriodically", arg0, arg1)
	ret0, _ := ret[0].(closeable.Closeable)
	return ret0
}

// SyncPeriodically indicates an expected call of SyncPeriodically.
func (mr *MockSyncMockRecorder) SyncPeriodically(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncPeriodically", reflect.TypeOf((*MockSync)(nil).SyncPeriodically), arg0, arg1)
}

// MockRPC is a mock of RPC interface.
type MockRPC struct {
	ctrl     *gomock.Controller
	recorder *MockRPCMockRecorder
}

// MockRPCMockRecorder is the mock recorder for MockRPC.
type MockRPCMockRecorder struct {
	mock *MockRPC
}

// NewMockRPC creates a new mock instance.
func NewMockRPC(ctrl *gomock.Controller) *MockRPC {
	mock := &MockRPC{ctrl: ctrl}
	mock.recorder = &MockRPCMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRPC) EXPECT() *MockRPCMockRecorder {
	return m.recorder
}

// GetClient mocks base method.
func (m *MockRPC) GetClient(arg0 entities.Peer) (proto.MetalInfraConfigClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClient", arg0)
	ret0, _ := ret[0].(proto.MetalInfraConfigClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClient indicates an expected call of GetClient.
func (mr *MockRPCMockRecorder) GetClient(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClient", reflect.TypeOf((*MockRPC)(nil).GetClient), arg0)
}
