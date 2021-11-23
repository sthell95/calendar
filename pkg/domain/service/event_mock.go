// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/domain/service/request_access.go

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	entity "calendar.com/pkg/domain/entity"
)

// MockEvent is a mock of Event interface.
type MockEvent struct {
	ctrl     *gomock.Controller
	recorder *MockEventMockRecorder
}

// MockEventMockRecorder is the mock recorder for MockEvent.
type MockEventMockRecorder struct {
	mock *MockEvent
}

// NewMockEvent creates a new mock instance.
func NewMockEvent(ctrl *gomock.Controller) *MockEvent {
	mock := &MockEvent{ctrl: ctrl}
	mock.recorder = &MockEventMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEvent) EXPECT() *MockEventMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockEvent) Create(event *entity.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", event)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockEventMockRecorder) Create(event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockEvent)(nil).Create), event)
}

// Update mocks base method.
func (m *MockEvent) Update(event *entity.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", event)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockEventMockRecorder) Update(event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockEvent)(nil).Update), event)
}
