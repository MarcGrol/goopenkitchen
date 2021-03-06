// Code generated by MockGen. DO NOT EDIT.
// Source: pincodeGenerator.go

// Package pincoder is a generated GoMock package.
package pincoder

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPincodeGenerator is a mock of PincodeGenerator interface.
type MockPincodeGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockPincodeGeneratorMockRecorder
}

// MockPincodeGeneratorMockRecorder is the mock recorder for MockPincodeGenerator.
type MockPincodeGeneratorMockRecorder struct {
	mock *MockPincodeGenerator
}

// NewMockPincodeGenerator creates a new mock instance.
func NewMockPincodeGenerator(ctrl *gomock.Controller) *MockPincodeGenerator {
	mock := &MockPincodeGenerator{ctrl: ctrl}
	mock.recorder = &MockPincodeGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPincodeGenerator) EXPECT() *MockPincodeGeneratorMockRecorder {
	return m.recorder
}

// GeneratePincode mocks base method.
func (m *MockPincodeGenerator) GeneratePincode() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GeneratePincode")
	ret0, _ := ret[0].(int)
	return ret0
}

// GeneratePincode indicates an expected call of GeneratePincode.
func (mr *MockPincodeGeneratorMockRecorder) GeneratePincode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GeneratePincode", reflect.TypeOf((*MockPincodeGenerator)(nil).GeneratePincode))
}
