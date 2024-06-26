// Code generated by MockGen. DO NOT EDIT.
// Source: bl/auth/authRepo/authRepo.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIAuthRepository is a mock of IAuthRepository interface.
type MockIAuthRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIAuthRepositoryMockRecorder
}

// MockIAuthRepositoryMockRecorder is the mock recorder for MockIAuthRepository.
type MockIAuthRepositoryMockRecorder struct {
	mock *MockIAuthRepository
}

// NewMockIAuthRepository creates a new mock instance.
func NewMockIAuthRepository(ctrl *gomock.Controller) *MockIAuthRepository {
	mock := &MockIAuthRepository{ctrl: ctrl}
	mock.recorder = &MockIAuthRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAuthRepository) EXPECT() *MockIAuthRepositoryMockRecorder {
	return m.recorder
}

// AddToken mocks base method.
func (m *MockIAuthRepository) AddToken(token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToken", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToken indicates an expected call of AddToken.
func (mr *MockIAuthRepositoryMockRecorder) AddToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToken", reflect.TypeOf((*MockIAuthRepository)(nil).AddToken), token)
}
