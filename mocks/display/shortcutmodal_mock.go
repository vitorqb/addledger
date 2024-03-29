// Code generated by MockGen. DO NOT EDIT.
// Source: shortcutmodal.go

// Package mock_display is a generated GoMock package.
package mock_display

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockShortcutModalController is a mock of ShortcutModalController interface.
type MockShortcutModalController struct {
	ctrl     *gomock.Controller
	recorder *MockShortcutModalControllerMockRecorder
}

// MockShortcutModalControllerMockRecorder is the mock recorder for MockShortcutModalController.
type MockShortcutModalControllerMockRecorder struct {
	mock *MockShortcutModalController
}

// NewMockShortcutModalController creates a new mock instance.
func NewMockShortcutModalController(ctrl *gomock.Controller) *MockShortcutModalController {
	mock := &MockShortcutModalController{ctrl: ctrl}
	mock.recorder = &MockShortcutModalControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShortcutModalController) EXPECT() *MockShortcutModalControllerMockRecorder {
	return m.recorder
}

// OnDiscardStatement mocks base method.
func (m *MockShortcutModalController) OnDiscardStatement() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnDiscardStatement")
}

// OnDiscardStatement indicates an expected call of OnDiscardStatement.
func (mr *MockShortcutModalControllerMockRecorder) OnDiscardStatement() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnDiscardStatement", reflect.TypeOf((*MockShortcutModalController)(nil).OnDiscardStatement))
}

// OnHideShortcutModal mocks base method.
func (m *MockShortcutModalController) OnHideShortcutModal() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnHideShortcutModal")
}

// OnHideShortcutModal indicates an expected call of OnHideShortcutModal.
func (mr *MockShortcutModalControllerMockRecorder) OnHideShortcutModal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnHideShortcutModal", reflect.TypeOf((*MockShortcutModalController)(nil).OnHideShortcutModal))
}

// OnLoadStatementRequest mocks base method.
func (m *MockShortcutModalController) OnLoadStatementRequest() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnLoadStatementRequest")
}

// OnLoadStatementRequest indicates an expected call of OnLoadStatementRequest.
func (mr *MockShortcutModalControllerMockRecorder) OnLoadStatementRequest() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnLoadStatementRequest", reflect.TypeOf((*MockShortcutModalController)(nil).OnLoadStatementRequest))
}
