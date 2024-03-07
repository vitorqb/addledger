// Code generated by MockGen. DO NOT EDIT.
// Source: app.go

// Package mock_context is a generated GoMock package.
package mock_context

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	tview "github.com/rivo/tview"
)

// MockTviewApp is a mock of TviewApp interface.
type MockTviewApp struct {
	ctrl     *gomock.Controller
	recorder *MockTviewAppMockRecorder
}

// MockTviewAppMockRecorder is the mock recorder for MockTviewApp.
type MockTviewAppMockRecorder struct {
	mock *MockTviewApp
}

// NewMockTviewApp creates a new mock instance.
func NewMockTviewApp(ctrl *gomock.Controller) *MockTviewApp {
	mock := &MockTviewApp{ctrl: ctrl}
	mock.recorder = &MockTviewAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTviewApp) EXPECT() *MockTviewAppMockRecorder {
	return m.recorder
}

// QueueUpdateDraw mocks base method.
func (m *MockTviewApp) QueueUpdateDraw(arg0 func()) *tview.Application {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueueUpdateDraw", arg0)
	ret0, _ := ret[0].(*tview.Application)
	return ret0
}

// QueueUpdateDraw indicates an expected call of QueueUpdateDraw.
func (mr *MockTviewAppMockRecorder) QueueUpdateDraw(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueUpdateDraw", reflect.TypeOf((*MockTviewApp)(nil).QueueUpdateDraw), arg0)
}