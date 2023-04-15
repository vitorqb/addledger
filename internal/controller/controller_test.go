package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	statemod "github.com/vitorqb/addledger/internal/state"
)

func TestInputController(t *testing.T) {
	t.Run("OnAccountInput and empty account", func(t *testing.T) {
		state := statemod.InitialState()
		controller := NewController(state)
		controller.OnAccountInput("")
		assert.Equal(t, statemod.Confirmation, state.CurrentPhase())
	})

	t.Run("OnAccountInput not empty", func(t *testing.T) {
		state := statemod.InitialState()
		state.SetPhase(statemod.InputPostingAccount)
		controller := NewController(state)
		controller.OnAccountInput("FOO")
		assert.Equal(t, statemod.InputPostingValue, state.CurrentPhase())
		account, _ := state.JournalEntryInput.CurrentPosting().GetAccount()
		assert.Equal(t, "FOO", account)
	})
}
