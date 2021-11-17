// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/domain/repository/event.go

// Package repository is a generated GoMock package.
package repository

import (
	reflect "reflect"

	entity "calendar.com/pkg/domain/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockEventRepository is a mock of EventRepository interface.
type MockEventRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEventRepositoryMockRecorder
}

// MockEventRepositoryMockRecorder is the mock recorder for MockEventRepository.
type MockEventRepositoryMockRecorder struct {
	mock *MockEventRepository
}

// NewMockEventRepository creates a new mock instance.
func NewMockEventRepository(ctrl *gomock.Controller) *MockEventRepository {
	mock := &MockEventRepository{ctrl: ctrl}
	mock.recorder = &MockEventRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventRepository) EXPECT() *MockEventRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockEventRepository) Create(event *entity.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", event)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockEventRepositoryMockRecorder) Create(event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockEventRepository)(nil).Create), event)
}

// FindOneById mocks base method.
func (m *MockEventRepository) FindOneById(arg0 string) (*entity.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneById", arg0)
	ret0, _ := ret[0].(*entity.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneById indicates an expected call of FindOneById.
func (mr *MockEventRepositoryMockRecorder) FindOneById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneById", reflect.TypeOf((*MockEventRepository)(nil).FindOneById), arg0)
}

// Update mocks base method.
func (m *MockEventRepository) Update(event *entity.Event, id string) (*entity.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", event, id)
	ret0, _ := ret[0].(*entity.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockEventRepositoryMockRecorder) Update(event, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockEventRepository)(nil).Update), event, id)
}
