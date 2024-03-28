// Code generated by MockGen. DO NOT EDIT.
// Source: bl/userService/userRepo/userRepository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	models "annotater/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIUserRepository is a mock of IUserRepository interface.
type MockIUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIUserRepositoryMockRecorder
}

// MockIUserRepositoryMockRecorder is the mock recorder for MockIUserRepository.
type MockIUserRepositoryMockRecorder struct {
	mock *MockIUserRepository
}

// NewMockIUserRepository creates a new mock instance.
func NewMockIUserRepository(ctrl *gomock.Controller) *MockIUserRepository {
	mock := &MockIUserRepository{ctrl: ctrl}
	mock.recorder = &MockIUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUserRepository) EXPECT() *MockIUserRepositoryMockRecorder {
	return m.recorder
}

// AddUser mocks base method.
func (m *MockIUserRepository) AddUser(user models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUser indicates an expected call of AddUser.
func (mr *MockIUserRepositoryMockRecorder) AddUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockIUserRepository)(nil).AddUser), user)
}

// DeleteUserByLogin mocks base method.
func (m *MockIUserRepository) DeleteUserByLogin(login string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserByLogin", login)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserByLogin indicates an expected call of DeleteUserByLogin.
func (mr *MockIUserRepositoryMockRecorder) DeleteUserByLogin(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserByLogin", reflect.TypeOf((*MockIUserRepository)(nil).DeleteUserByLogin), login)
}

// GetUserByCookie mocks base method.
func (m *MockIUserRepository) GetUserByCookie(cookie models.Cookie) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByCookie", cookie)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByCookie indicates an expected call of GetUserByCookie.
func (mr *MockIUserRepositoryMockRecorder) GetUserByCookie(cookie interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByCookie", reflect.TypeOf((*MockIUserRepository)(nil).GetUserByCookie), cookie)
}

// GetUserByID mocks base method.
func (m *MockIUserRepository) GetUserByID(id uint64) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", id)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockIUserRepositoryMockRecorder) GetUserByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockIUserRepository)(nil).GetUserByID), id)
}

// GetUserByLogin mocks base method.
func (m *MockIUserRepository) GetUserByLogin(login string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", login)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockIUserRepositoryMockRecorder) GetUserByLogin(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockIUserRepository)(nil).GetUserByLogin), login)
}

// UpdateUserByLogin mocks base method.
func (m *MockIUserRepository) UpdateUserByLogin(login string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserByLogin", login)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserByLogin indicates an expected call of UpdateUserByLogin.
func (mr *MockIUserRepositoryMockRecorder) UpdateUserByLogin(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserByLogin", reflect.TypeOf((*MockIUserRepository)(nil).UpdateUserByLogin), login)
}
