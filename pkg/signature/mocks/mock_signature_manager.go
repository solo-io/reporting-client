// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/solo-io/reporting-client/pkg/sig (interfaces: SignatureManager)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSignatureManager is a mock of SignatureManager interface
type MockSignatureManager struct {
	ctrl     *gomock.Controller
	recorder *MockSignatureManagerMockRecorder
}

// MockSignatureManagerMockRecorder is the mock recorder for MockSignatureManager
type MockSignatureManagerMockRecorder struct {
	mock *MockSignatureManager
}

// NewMockSignatureManager creates a new mock instance
func NewMockSignatureManager(ctrl *gomock.Controller) *MockSignatureManager {
	mock := &MockSignatureManager{ctrl: ctrl}
	mock.recorder = &MockSignatureManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSignatureManager) EXPECT() *MockSignatureManagerMockRecorder {
	return m.recorder
}

// GetSignature mocks base method
func (m *MockSignatureManager) GetSignature() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSignature")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSignature indicates an expected call of GetSignature
func (mr *MockSignatureManagerMockRecorder) GetSignature() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSignature", reflect.TypeOf((*MockSignatureManager)(nil).GetSignature))
}