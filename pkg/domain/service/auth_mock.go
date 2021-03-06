// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/domain/service/auth.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	entity "calendar.com/pkg/domain/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockCredentials is a mock of Credentials interface.
type MockCredentials struct {
	ctrl     *gomock.Controller
	recorder *MockCredentialsMockRecorder
}

// MockCredentialsMockRecorder is the mock recorder for MockCredentials.
type MockCredentialsMockRecorder struct {
	mock *MockCredentials
}

// NewMockCredentials creates a new mock instance.
func NewMockCredentials(ctrl *gomock.Controller) *MockCredentials {
	mock := &MockCredentials{ctrl: ctrl}
	mock.recorder = &MockCredentialsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCredentials) EXPECT() *MockCredentialsMockRecorder {
	return m.recorder
}

// CheckCredentials mocks base method.
func (m *MockCredentials) CheckCredentials(arg0 entity.Credentials) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckCredentials", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckCredentials indicates an expected call of CheckCredentials.
func (mr *MockCredentialsMockRecorder) CheckCredentials(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckCredentials", reflect.TypeOf((*MockCredentials)(nil).CheckCredentials), arg0)
}

// GenerateToken mocks base method.
func (m *MockCredentials) GenerateToken(arg0 *entity.Credentials) (*entity.AuthToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", arg0)
	ret0, _ := ret[0].(*entity.AuthToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockCredentialsMockRecorder) GenerateToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockCredentials)(nil).GenerateToken), arg0)
}

// MockAuthorization is a mock of Authorization interface.
type MockAuthorization struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationMockRecorder
}

// MockAuthorizationMockRecorder is the mock recorder for MockAuthorization.
type MockAuthorizationMockRecorder struct {
	mock *MockAuthorization
}

// NewMockAuthorization creates a new mock instance.
func NewMockAuthorization(ctrl *gomock.Controller) *MockAuthorization {
	mock := &MockAuthorization{ctrl: ctrl}
	mock.recorder = &MockAuthorizationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorization) EXPECT() *MockAuthorizationMockRecorder {
	return m.recorder
}

// SignInProcess mocks base method.
func (m *MockAuthorization) SignInProcess(ctx context.Context, c *entity.Credentials) (*entity.AuthToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignInProcess", ctx, c)
	ret0, _ := ret[0].(*entity.AuthToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignInProcess indicates an expected call of SignInProcess.
func (mr *MockAuthorizationMockRecorder) SignInProcess(ctx, c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignInProcess", reflect.TypeOf((*MockAuthorization)(nil).SignInProcess), ctx, c)
}
