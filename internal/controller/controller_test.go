package controller

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	statemod "github.com/vitorqb/addledger/internal/state"
)

func TestInputController(t *testing.T) {
	t.Run("OnDateInput", func(t *testing.T) {
		state := statemod.InitialState()
		state.SetPhase(statemod.InputDate)
		controller := NewController(state)
		date, _ := time.Parse(time.RFC3339, "2022-01-01")
		controller.OnDateInput(date)
		assert.Equal(t, statemod.InputDescription, state.CurrentPhase())
		foundDate, _ := state.JournalEntryInput.GetDate()
		assert.Equal(t, date, foundDate)
	})
	t.Run("OnDescriptionInput", func(t *testing.T) {
		state := statemod.InitialState()
		state.SetPhase(statemod.InputDescription)
		controller := NewController(state)
		controller.OnDescriptionInput("FOO")
		assert.Equal(t, statemod.InputPostingAccount, state.CurrentPhase())
		foundDescription, _ := state.JournalEntryInput.GetDescription()
		assert.Equal(t, "FOO", foundDescription)
	})
	t.Run("OnAccountInput and empty account", func(t *testing.T) {
		state := statemod.InitialState()
		controller := NewController(state)
		controller.OnPostingAccountInput("")
		assert.Equal(t, statemod.Confirmation, state.CurrentPhase())
	})

	t.Run("OnAccountInput not empty", func(t *testing.T) {
		state := statemod.InitialState()
		state.SetPhase(statemod.InputPostingAccount)
		controller := NewController(state)
		controller.OnPostingAccountInput("FOO")
		assert.Equal(t, statemod.InputPostingValue, state.CurrentPhase())
		account, _ := state.JournalEntryInput.CurrentPosting().GetAccount()
		assert.Equal(t, "FOO", account)
	})
}
