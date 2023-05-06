package controller

import (
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/internal/listaction"
	statemod "github.com/vitorqb/addledger/internal/state"
)

//go:generate mockgen --source=controller.go --destination=../../mocks/controller/controller_mock.go

// IInputController reacts to the user inputs and interactions.
type IInputController interface {
	OnDateInput(date time.Time)
	OnDescriptionInput(description string)
	OnPostingAccountDone(account string)
	OnPostingAccountListAcction(action listaction.ListAction)
	OnPostingValueInput(value string)
	OnInputConfirmation()
	OnInputRejection()
}

// InputController implements IInputController.
type InputController struct {
	state    *statemod.State
	output   io.Writer
	eventBus eventbus.IEventBus
}

var _ IInputController = &InputController{}

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

func (ic *InputController) OnPostingAccountListAcction(action listaction.ListAction) {
	err := ic.eventBus.Send(eventbus.Event{
		Topic: "input.postingaccount.listaction",
		Data:  action,
	})
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
