// Code generated by MockGen. DO NOT EDIT.
// Source: hledger.go

// Package mock_hledger is a generated GoMock package.
package mock_hledger

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	journal "github.com/vitorqb/addledger/internal/journal"
)

// MockIClient is a mock of IClient interface.
type MockIClient struct {
	ctrl     *gomock.Controller
	recorder *MockIClientMockRecorder
}

// MockIClientMockRecorder is the mock recorder for MockIClient.
type MockIClientMockRecorder struct {
	mock *MockIClient
}

// NewMockIClient creates a new mock instance.
func NewMockIClient(ctrl *gomock.Controller) *MockIClient {
	mock := &MockIClient{ctrl: ctrl}
	mock.recorder = &MockIClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIClient) EXPECT() *MockIClientMockRecorder {
	return m.recorder
}

// Accounts mocks base method.
func (m *MockIClient) Accounts() ([]journal.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accounts")
	ret0, _ := ret[0].([]journal.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Accounts indicates an expected call of Accounts.
func (mr *MockIClientMockRecorder) Accounts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accounts", reflect.TypeOf((*MockIClient)(nil).Accounts))
}

// Transactions mocks base method.
func (m *MockIClient) Transactions() ([]journal.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transactions")
	ret0, _ := ret[0].([]journal.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Transactions indicates an expected call of Transactions.
func (mr *MockIClientMockRecorder) Transactions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transactions", reflect.TypeOf((*MockIClient)(nil).Transactions))
}
