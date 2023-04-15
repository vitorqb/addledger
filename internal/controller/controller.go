package controller

import (
	"time"

	statemod "github.com/vitorqb/addledger/internal/state"
)

type (
	InputController struct {
		state *statemod.State
	}
)

func NewController(state *statemod.State) *InputController {
	return &InputController{state}
}

func (ic *InputController) OnDateInput(date time.Time) {
	ic.state.JournalEntryInput.SetDate(date)
	ic.state.NextPhase()
}

func (ic *InputController) OnDescriptionInput(description string) {
	ic.state.JournalEntryInput.SetDescription(description)
	ic.state.NextPhase()
}

func (ic *InputController) OnPostingAccountInput(account string) {
	// Empty string -> user is done entering postings.
	if account == "" {
		ic.state.SetPhase(statemod.Confirmation)
		return
	}

	// Otherwise save the account and move to value
	posting := ic.state.JournalEntryInput.CurrentPosting()
	posting.SetAccount(account)
	ic.state.NextPhase()
}

func (ic *InputController) OnPostingValueInput(value string) {
	posting := ic.state.JournalEntryInput.CurrentPosting()
	posting.SetValue(value)
	ic.state.JournalEntryInput.AdvancePosting()
	ic.state.SetPhase(statemod.InputPostingAccount)
}
