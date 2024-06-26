// Code generated by MockGen. DO NOT EDIT.
// Source: bl/auth/authService.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	models "annotater/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIAuthService is a mock of IAuthService interface.
type MockIAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockIAuthServiceMockRecorder
}

// MockIAuthServiceMockRecorder is the mock recorder for MockIAuthService.
type MockIAuthServiceMockRecorder struct {
	mock *MockIAuthService
}

// NewMockIAuthService creates a new mock instance.
func NewMockIAuthService(ctrl *gomock.Controller) *MockIAuthService {
	mock := &MockIAuthService{ctrl: ctrl}
	mock.recorder = &MockIAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAuthService) EXPECT() *MockIAuthServiceMockRecorder {
	return m.recorder
}

// SignIn mocks base method.
func (m *MockIAuthService) SignIn(candidate *models.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", candidate)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockIAuthServiceMockRecorder) SignIn(candidate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockIAuthService)(nil).SignIn), candidate)
}

// SignUp mocks base method.
func (m *MockIAuthService) SignUp(candidate *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", candidate)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignUp indicates an expected call of SignUp.
func (mr *MockIAuthServiceMockRecorder) SignUp(candidate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockIAuthService)(nil).SignUp), candidate)
}
