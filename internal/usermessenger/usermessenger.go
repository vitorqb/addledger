package usermessenger

import (
	"fmt"

	statemod "github.com/vitorqb/addledger/internal/state"
)

//go:generate $MOCKGEN --source=usermessenger.go --destination=../../mocks/usermessenger/usermessenger_mock.go

type IUserMessenger interface {
	Info(string)
	Warning(msg string, err error)
	Error(msg string, err error)
}

type UserMessenger struct {
	state *statemod.State
}

func (u *UserMessenger) Info(msg string) {
	u.state.Display.SetUserMessage(msg)
}

func (u *UserMessenger) Warning(msg string, err error) {
	warnMsg := "WARNING: " + msg
	if err != nil {
		warnMsg += fmt.Sprintf(" - %s", err.Error())
	}
	u.state.Display.SetUserMessage(warnMsg)
}

func (u *UserMessenger) Error(msg string, err error) {
	errMsg := "ERROR: " + msg
	if err != nil {
		errMsg += fmt.Sprintf(" - %s", err.Error())
	}
	u.state.Display.SetUserMessage(errMsg)
}

func New(state *statemod.State) *UserMessenger {
	return &UserMessenger{state}
}

type NoOp struct{}

func (n *NoOp) Info(string)           {}
func (n *NoOp) Warning(string, error) {}
func (n *NoOp) Error(string, error)   {}
