package controller

import statemod "github.com/vitorqb/addledger/internal/state"

type (
	InputController struct {
		state *statemod.State
	}
)

func NewController(state *statemod.State) *InputController {
	return &InputController{state}
}

// !!! TODO Rename to OnPostingAccountInput
func (ic *InputController) OnAccountInput(account string) {
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
