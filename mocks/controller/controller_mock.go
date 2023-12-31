// Code generated by MockGen. DO NOT EDIT.
// Source: controller.go

// Package mock_controller is a generated GoMock package.
package mock_controller

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	config "github.com/vitorqb/addledger/internal/config"
	listaction "github.com/vitorqb/addledger/internal/listaction"
	userinput "github.com/vitorqb/addledger/internal/userinput"
)

// MockStatementLoader is a mock of StatementLoader interface.
type MockStatementLoader struct {
	ctrl     *gomock.Controller
	recorder *MockStatementLoaderMockRecorder
}

// MockStatementLoaderMockRecorder is the mock recorder for MockStatementLoader.
type MockStatementLoaderMockRecorder struct {
	mock *MockStatementLoader
}

// NewMockStatementLoader creates a new mock instance.
func NewMockStatementLoader(ctrl *gomock.Controller) *MockStatementLoader {
	mock := &MockStatementLoader{ctrl: ctrl}
	mock.recorder = &MockStatementLoaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStatementLoader) EXPECT() *MockStatementLoaderMockRecorder {
	return m.recorder
}

// Load mocks base method.
func (m *MockStatementLoader) Load(config config.StatementLoaderConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", config)
	ret0, _ := ret[0].(error)
	return ret0
}

// Load indicates an expected call of Load.
func (mr *MockStatementLoaderMockRecorder) Load(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockStatementLoader)(nil).Load), config)
}

// MockIInputController is a mock of IInputController interface.
type MockIInputController struct {
	ctrl     *gomock.Controller
	recorder *MockIInputControllerMockRecorder
}

// MockIInputControllerMockRecorder is the mock recorder for MockIInputController.
type MockIInputControllerMockRecorder struct {
	mock *MockIInputController
}

// NewMockIInputController creates a new mock instance.
func NewMockIInputController(ctrl *gomock.Controller) *MockIInputController {
	mock := &MockIInputController{ctrl: ctrl}
	mock.recorder = &MockIInputControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIInputController) EXPECT() *MockIInputControllerMockRecorder {
	return m.recorder
}

// OnDateChanged mocks base method.
func (m *MockIInputController) OnDateChanged(text string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnDateChanged", text)
}

// OnDateChanged indicates an expected call of OnDateChanged.
func (mr *MockIInputControllerMockRecorder) OnDateChanged(text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnDateChanged", reflect.TypeOf((*MockIInputController)(nil).OnDateChanged), text)
}

// OnDateDone mocks base method.
func (m *MockIInputController) OnDateDone() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnDateDone")
}

// OnDateDone indicates an expected call of OnDateDone.
func (mr *MockIInputControllerMockRecorder) OnDateDone() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnDateDone", reflect.TypeOf((*MockIInputController)(nil).OnDateDone))
}

// OnDescriptionChanged mocks base method.
func (m *MockIInputController) OnDescriptionChanged(newText string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnDescriptionChanged", newText)
}

// OnDescriptionChanged indicates an expected call of OnDescriptionChanged.
func (mr *MockIInputControllerMockRecorder) OnDescriptionChanged(newText interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnDescriptionChanged", reflect.TypeOf((*MockIInputController)(nil).OnDescriptionChanged), newText)
}

// OnDescriptionDone mocks base method.
func (m *MockIInputController) OnDescriptionDone(source userinput.DoneSource) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnDescriptionDone", source)
}

// OnDescriptionDone indicates an expected call of OnDescriptionDone.
func (mr *MockIInputControllerMockRecorder) OnDescriptionDone(source interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnDescriptionDone", reflect.TypeOf((*MockIInputController)(nil).OnDescriptionDone), source)
}

// OnDescriptionInsertFromContext mocks base method.
func (m *MockIInputController) OnDescriptionInsertFromContext() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnDescriptionInsertFromContext")
}

// OnDescriptionInsertFromContext indicates an expected call of OnDescriptionInsertFromContext.
func (mr *MockIInputControllerMockRecorder) OnDescriptionInsertFromContext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnDescriptionInsertFromContext", reflect.TypeOf((*MockIInputController)(nil).OnDescriptionInsertFromContext))
}

// OnDescriptionListAction mocks base method.
func (m *MockIInputController) OnDescriptionListAction(action listaction.ListAction) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnDescriptionListAction", action)
}

// OnDescriptionListAction indicates an expected call of OnDescriptionListAction.
func (mr *MockIInputControllerMockRecorder) OnDescriptionListAction(action interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnDescriptionListAction", reflect.TypeOf((*MockIInputController)(nil).OnDescriptionListAction), action)
}

// OnDiscardStatement mocks base method.
func (m *MockIInputController) OnDiscardStatement() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnDiscardStatement")
}

// OnDiscardStatement indicates an expected call of OnDiscardStatement.
func (mr *MockIInputControllerMockRecorder) OnDiscardStatement() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnDiscardStatement", reflect.TypeOf((*MockIInputController)(nil).OnDiscardStatement))
}

// OnDisplayShortcutModal mocks base method.
func (m *MockIInputController) OnDisplayShortcutModal() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnDisplayShortcutModal")
}

// OnDisplayShortcutModal indicates an expected call of OnDisplayShortcutModal.
func (mr *MockIInputControllerMockRecorder) OnDisplayShortcutModal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnDisplayShortcutModal", reflect.TypeOf((*MockIInputController)(nil).OnDisplayShortcutModal))
}

// OnFinishPosting mocks base method.
func (m *MockIInputController) OnFinishPosting() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnFinishPosting")
}

// OnFinishPosting indicates an expected call of OnFinishPosting.
func (mr *MockIInputControllerMockRecorder) OnFinishPosting() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnFinishPosting", reflect.TypeOf((*MockIInputController)(nil).OnFinishPosting))
}

// OnHideShortcutModal mocks base method.
func (m *MockIInputController) OnHideShortcutModal() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnHideShortcutModal")
}

// OnHideShortcutModal indicates an expected call of OnHideShortcutModal.
func (mr *MockIInputControllerMockRecorder) OnHideShortcutModal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnHideShortcutModal", reflect.TypeOf((*MockIInputController)(nil).OnHideShortcutModal))
}

// OnInputConfirmation mocks base method.
func (m *MockIInputController) OnInputConfirmation() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnInputConfirmation")
}

// OnInputConfirmation indicates an expected call of OnInputConfirmation.
func (mr *MockIInputControllerMockRecorder) OnInputConfirmation() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnInputConfirmation", reflect.TypeOf((*MockIInputController)(nil).OnInputConfirmation))
}

// OnInputRejection mocks base method.
func (m *MockIInputController) OnInputRejection() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnInputRejection")
}

// OnInputRejection indicates an expected call of OnInputRejection.
func (mr *MockIInputControllerMockRecorder) OnInputRejection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnInputRejection", reflect.TypeOf((*MockIInputController)(nil).OnInputRejection))
}

// OnLoadStatement mocks base method.
func (m *MockIInputController) OnLoadStatement(csvFile, presetFile string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnLoadStatement", csvFile, presetFile)
}

// OnLoadStatement indicates an expected call of OnLoadStatement.
func (mr *MockIInputControllerMockRecorder) OnLoadStatement(csvFile, presetFile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnLoadStatement", reflect.TypeOf((*MockIInputController)(nil).OnLoadStatement), csvFile, presetFile)
}

// OnLoadStatementRequest mocks base method.
func (m *MockIInputController) OnLoadStatementRequest() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnLoadStatementRequest")
}

// OnLoadStatementRequest indicates an expected call of OnLoadStatementRequest.
func (mr *MockIInputControllerMockRecorder) OnLoadStatementRequest() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnLoadStatementRequest", reflect.TypeOf((*MockIInputController)(nil).OnLoadStatementRequest))
}

// OnPostingAccountChanged mocks base method.
func (m *MockIInputController) OnPostingAccountChanged(newText string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnPostingAccountChanged", newText)
}

// OnPostingAccountChanged indicates an expected call of OnPostingAccountChanged.
func (mr *MockIInputControllerMockRecorder) OnPostingAccountChanged(newText interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnPostingAccountChanged", reflect.TypeOf((*MockIInputController)(nil).OnPostingAccountChanged), newText)
}

// OnPostingAccountDone mocks base method.
func (m *MockIInputController) OnPostingAccountDone(source userinput.DoneSource) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnPostingAccountDone", source)
}

// OnPostingAccountDone indicates an expected call of OnPostingAccountDone.
func (mr *MockIInputControllerMockRecorder) OnPostingAccountDone(source interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnPostingAccountDone", reflect.TypeOf((*MockIInputController)(nil).OnPostingAccountDone), source)
}

// OnPostingAccountInsertFromContext mocks base method.
func (m *MockIInputController) OnPostingAccountInsertFromContext() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnPostingAccountInsertFromContext")
}

// OnPostingAccountInsertFromContext indicates an expected call of OnPostingAccountInsertFromContext.
func (mr *MockIInputControllerMockRecorder) OnPostingAccountInsertFromContext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnPostingAccountInsertFromContext", reflect.TypeOf((*MockIInputController)(nil).OnPostingAccountInsertFromContext))
}

// OnPostingAccountListAcction mocks base method.
func (m *MockIInputController) OnPostingAccountListAcction(action listaction.ListAction) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnPostingAccountListAcction", action)
}

// OnPostingAccountListAcction indicates an expected call of OnPostingAccountListAcction.
func (mr *MockIInputControllerMockRecorder) OnPostingAccountListAcction(action interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnPostingAccountListAcction", reflect.TypeOf((*MockIInputController)(nil).OnPostingAccountListAcction), action)
}

// OnPostingAmmountChanged mocks base method.
func (m *MockIInputController) OnPostingAmmountChanged(text string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnPostingAmmountChanged", text)
}

// OnPostingAmmountChanged indicates an expected call of OnPostingAmmountChanged.
func (mr *MockIInputControllerMockRecorder) OnPostingAmmountChanged(text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnPostingAmmountChanged", reflect.TypeOf((*MockIInputController)(nil).OnPostingAmmountChanged), text)
}

// OnPostingAmmountDone mocks base method.
func (m *MockIInputController) OnPostingAmmountDone(arg0 userinput.DoneSource) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnPostingAmmountDone", arg0)
}

// OnPostingAmmountDone indicates an expected call of OnPostingAmmountDone.
func (mr *MockIInputControllerMockRecorder) OnPostingAmmountDone(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnPostingAmmountDone", reflect.TypeOf((*MockIInputController)(nil).OnPostingAmmountDone), arg0)
}

// OnTagChanged mocks base method.
func (m *MockIInputController) OnTagChanged(newText string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnTagChanged", newText)
}

// OnTagChanged indicates an expected call of OnTagChanged.
func (mr *MockIInputControllerMockRecorder) OnTagChanged(newText interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnTagChanged", reflect.TypeOf((*MockIInputController)(nil).OnTagChanged), newText)
}

// OnTagDone mocks base method.
func (m *MockIInputController) OnTagDone(source userinput.DoneSource) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnTagDone", source)
}

// OnTagDone indicates an expected call of OnTagDone.
func (mr *MockIInputControllerMockRecorder) OnTagDone(source interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnTagDone", reflect.TypeOf((*MockIInputController)(nil).OnTagDone), source)
}

// OnTagInsertFromContext mocks base method.
func (m *MockIInputController) OnTagInsertFromContext() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnTagInsertFromContext")
}

// OnTagInsertFromContext indicates an expected call of OnTagInsertFromContext.
func (mr *MockIInputControllerMockRecorder) OnTagInsertFromContext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnTagInsertFromContext", reflect.TypeOf((*MockIInputController)(nil).OnTagInsertFromContext))
}

// OnTagListAction mocks base method.
func (m *MockIInputController) OnTagListAction(action listaction.ListAction) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnTagListAction", action)
}

// OnTagListAction indicates an expected call of OnTagListAction.
func (mr *MockIInputControllerMockRecorder) OnTagListAction(action interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnTagListAction", reflect.TypeOf((*MockIInputController)(nil).OnTagListAction), action)
}

// OnUndo mocks base method.
func (m *MockIInputController) OnUndo() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnUndo")
}

// OnUndo indicates an expected call of OnUndo.
func (mr *MockIInputControllerMockRecorder) OnUndo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnUndo", reflect.TypeOf((*MockIInputController)(nil).OnUndo))
}
