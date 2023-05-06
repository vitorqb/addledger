package controller

import (
	"fmt"
	"io"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/input"
	statemod "github.com/vitorqb/addledger/internal/state"
)

type InputController struct {
	state    *statemod.State
	output   io.Writer
	eventBus eventbus.IEventBus
}

func NewController(state *statemod.State, options ...Opt) (*InputController, error) {
	opts := &Opts{}
	for _, option := range options {
		err := option.configure(opts)
		if err != nil {
			return nil, err
		}
	}
	if opts.output == nil {
		return nil, fmt.Errorf("missing output")
	}
	if opts.eventBus == nil {
		return nil, fmt.Errorf("missing Event Bus")
	}
	return &InputController{
		state:    state,
		output:   opts.output,
		eventBus: opts.eventBus,
	}, nil
}

func (ic *InputController) OnDateInput(date time.Time) {
	ic.state.JournalEntryInput.SetDate(date)
	ic.state.NextPhase()
}

func (ic *InputController) OnDescriptionInput(description string) {
	ic.state.JournalEntryInput.SetDescription(description)
	ic.state.NextPhase()
}

func (ic *InputController) OnPostingAccountDone(account string) {
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

func (ic *InputController) OnPostingAccountInputCapture(keyEvent *tcell.EventKey) {
	busEvent := eventbus.Event{
		Topic: "input.postingaccount.eventkey",
		Data:  keyEvent,
	}
	err := ic.eventBus.Send(busEvent)
	if err != nil {
		logrus.WithError(err).Warn("Failed to send event")
	}

}

func (ic *InputController) OnPostingValueInput(value string) {
	posting := ic.state.JournalEntryInput.CurrentPosting()
	posting.SetValue(value)
	ic.state.JournalEntryInput.AdvancePosting()
	ic.state.SetPhase(statemod.InputPostingAccount)
}

func (ic *InputController) OnInputConfirmation() {
	_, err := io.WriteString(ic.output, "\n\n"+ic.state.JournalEntryInput.Repr())
	if err != nil {
		// TODO Let user know somehow!
		logrus.WithError(err).Warn("failed to write to file")
		return
	}
	ic.state.JournalEntryInput = input.NewJournalEntryInput()
	ic.state.SetPhase(statemod.InputDate)
}

func (ic *InputController) OnInputRejection() {
	ic.state.SetPhase(statemod.InputPostingAccount)
}
