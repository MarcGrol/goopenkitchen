// Code generated by MockGen. DO NOT EDIT.
// Source: fileTimeWriter.go

// Package testable is a generated GoMock package.
package testable

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockUuidGenerator is a mock of UuidGenerator interface.
type MockUuidGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockUuidGeneratorMockRecorder
}

// MockUuidGeneratorMockRecorder is the mock recorder for MockUuidGenerator.
type MockUuidGeneratorMockRecorder struct {
	mock *MockUuidGenerator
}

// NewMockUuidGenerator creates a new mock instance.
func NewMockUuidGenerator(ctrl *gomock.Controller) *MockUuidGenerator {
	mock := &MockUuidGenerator{ctrl: ctrl}
	mock.recorder = &MockUuidGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUuidGenerator) EXPECT() *MockUuidGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockUuidGenerator) Generate() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate")
	ret0, _ := ret[0].(string)
	return ret0
}

// Generate indicates an expected call of Generate.
func (mr *MockUuidGeneratorMockRecorder) Generate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockUuidGenerator)(nil).Generate))
}

// MockNower is a mock of Nower interface.
type MockNower struct {
	ctrl     *gomock.Controller
	recorder *MockNowerMockRecorder
}

// MockNowerMockRecorder is the mock recorder for MockNower.
type MockNowerMockRecorder struct {
	mock *MockNower
}

// NewMockNower creates a new mock instance.
func NewMockNower(ctrl *gomock.Controller) *MockNower {
	mock := &MockNower{ctrl: ctrl}
	mock.recorder = &MockNowerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNower) EXPECT() *MockNowerMockRecorder {
	return m.recorder
}

// Now mocks base method.
func (m *MockNower) Now() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Now indicates an expected call of Now.
func (mr *MockNowerMockRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*MockNower)(nil).Now))
}
