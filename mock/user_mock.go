// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/impactasaurus/server/auth (interfaces: User)

// Package mock_auth is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUser is a mock of User interface
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// GetAssessmentScope mocks base method
func (m *MockUser) GetAssessmentScope() (string, bool) {
	ret := m.ctrl.Call(m, "GetAssessmentScope")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetAssessmentScope indicates an expected call of GetAssessmentScope
func (mr *MockUserMockRecorder) GetAssessmentScope() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAssessmentScope", reflect.TypeOf((*MockUser)(nil).GetAssessmentScope))
}

// IsBeneficiary mocks base method
func (m *MockUser) IsBeneficiary() bool {
	ret := m.ctrl.Call(m, "IsBeneficiary")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsBeneficiary indicates an expected call of IsBeneficiary
func (mr *MockUserMockRecorder) IsBeneficiary() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsBeneficiary", reflect.TypeOf((*MockUser)(nil).IsBeneficiary))
}

// Organisation mocks base method
func (m *MockUser) Organisation() (string, error) {
	ret := m.ctrl.Call(m, "Organisation")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Organisation indicates an expected call of Organisation
func (mr *MockUserMockRecorder) Organisation() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Organisation", reflect.TypeOf((*MockUser)(nil).Organisation))
}

// UserID mocks base method
func (m *MockUser) UserID() string {
	ret := m.ctrl.Call(m, "UserID")
	ret0, _ := ret[0].(string)
	return ret0
}

// UserID indicates an expected call of UserID
func (mr *MockUserMockRecorder) UserID() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserID", reflect.TypeOf((*MockUser)(nil).UserID))
}
