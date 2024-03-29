// Code generated by MockGen. DO NOT EDIT.
// Source: ammountguesser.go

// Package mock_ammountguesser is a generated GoMock package.
package mock_ammountguesser

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	ammountguesser "github.com/vitorqb/addledger/internal/ammountguesser"
	finance "github.com/vitorqb/addledger/internal/finance"
)

// MockIAmmountGuesser is a mock of IAmmountGuesser interface.
type MockIAmmountGuesser struct {
	ctrl     *gomock.Controller
	recorder *MockIAmmountGuesserMockRecorder
}

// MockIAmmountGuesserMockRecorder is the mock recorder for MockIAmmountGuesser.
type MockIAmmountGuesserMockRecorder struct {
	mock *MockIAmmountGuesser
}

// NewMockIAmmountGuesser creates a new mock instance.
func NewMockIAmmountGuesser(ctrl *gomock.Controller) *MockIAmmountGuesser {
	mock := &MockIAmmountGuesser{ctrl: ctrl}
	mock.recorder = &MockIAmmountGuesserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAmmountGuesser) EXPECT() *MockIAmmountGuesserMockRecorder {
	return m.recorder
}

// Guess mocks base method.
func (m *MockIAmmountGuesser) Guess(inputs ammountguesser.Inputs) (finance.Ammount, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Guess", inputs)
	ret0, _ := ret[0].(finance.Ammount)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Guess indicates an expected call of Guess.
func (mr *MockIAmmountGuesserMockRecorder) Guess(inputs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Guess", reflect.TypeOf((*MockIAmmountGuesser)(nil).Guess), inputs)
}
