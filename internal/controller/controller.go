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
	OnPostingValueInput(value string)
	OnInputConfirmation()
	OnInputRejection()

	// Called when an user wants to undo it's last action.
	OnUndo()

	// Controls Posting Account input
	OnPostingAccountChanged(newText string)
	OnPostingAccountDone(account string)
	OnPostingAccountInsertFromContext()
	OnPostingAccountListAcction(action listaction.ListAction)
	OnPostingAccountSelectedFromContext()

	// Controls the Description input
	OnDescriptionChanged(newText string)
	OnDescriptionDone()
	OnDescriptionInsertFromContext()
	OnDescriptionListAction(action listaction.ListAction)
	OnDescriptionSelectedFromContext()
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

func (ic *InputController) OnPostingAccountDone(account string) {
	// Empty string -> user is done entering postings.
	if account == "" {
		// remove the current posting since it's empty
		ic.state.JournalEntryInput.DeleteCurrentPosting()
		ic.state.SetPhase(statemod.Confirmation)
		return
	}

	// Otherwise save the account and move to value
	posting := ic.state.JournalEntryInput.CurrentPosting()
	posting.SetAccount(account)
	ic.state.NextPhase()
}

func (ic *InputController) OnPostingAccountSelectedFromContext() {
	selectedAccountFromContext := ic.state.InputMetadata.SelectedPostingAccount()
	ic.OnPostingAccountDone(selectedAccountFromContext)
}

// OnPostingAccountInsertFromContext inserts the text from the context to the
// PostingAccount input.
func (ic *InputController) OnPostingAccountInsertFromContext() {
	textFromContext := ic.state.InputMetadata.SelectedPostingAccount()
	event := eventbus.Event{
		Topic: "input.postingaccount.settext",
		Data:  textFromContext,
	}
	err := ic.eventBus.Send(event)
	if err != nil {
		logrus.WithError(err).Warn("Failed to send event")
	}
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

func (ic *InputController) OnPostingAccountChanged(newText string) {
	ic.state.InputMetadata.SetPostingAccountText(newText)
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
		logrus.WithError(err).Fatal("failed to write to file")
		return
	}
	ic.state.JournalEntryInput = input.NewJournalEntryInput()
	ic.state.SetPhase(statemod.InputDate)
}

func (ic *InputController) OnInputRejection() {
	// put back an empty posting so the user can add to it
	ic.state.JournalEntryInput.AdvancePosting()
	ic.state.SetPhase(statemod.InputPostingAccount)
}

func (ic *InputController) OnDescriptionChanged(newText string) {
	ic.state.InputMetadata.SetDescriptionText(newText)
}

func (ic *InputController) OnDescriptionListAction(action listaction.ListAction) {
	err := ic.eventBus.Send(eventbus.Event{
		Topic: "input.description.listaction",
		Data:  action,
	})
	if err != nil {
		logrus.WithError(err).Warn("Failed to send event")
	}
}

func (ic *InputController) OnDescriptionSelectedFromContext() {
	descriptionFromContext := ic.state.InputMetadata.SelectedDescription()
	ic.OnDescriptionChanged(descriptionFromContext)
	ic.OnDescriptionDone()
}

func (ic *InputController) OnDescriptionDone() {
	description := ic.state.InputMetadata.DescriptionText()
	ic.state.JournalEntryInput.SetDescription(description)
	if ic.state.JournalEntryInput.CountPostings() == 0 {
		ic.state.JournalEntryInput.AddPosting()
	}
	ic.state.NextPhase()
}

func (ic *InputController) OnDescriptionInsertFromContext() {
	descriptionFromContext := ic.state.InputMetadata.SelectedDescription()
	event := eventbus.Event{
		Topic: "input.description.settext",
		Data:  descriptionFromContext,
	}
	err := ic.eventBus.Send(event)
	if err != nil {
		logrus.WithError(err).Warn("Failed to send event")
	}
}

func (ic *InputController) OnUndo() {
	switch ic.state.CurrentPhase() {
	case statemod.InputDate:
		ic.state.PrevPhase()
	case statemod.InputDescription:
		ic.state.JournalEntryInput.ClearDate()
		ic.state.PrevPhase()
	case statemod.InputPostingAccount:
		ic.state.JournalEntryInput.DeleteCurrentPosting()
		if ic.state.JournalEntryInput.CountPostings() == 0 {
			// We don't have any postings - clear description and go back
			ic.state.JournalEntryInput.ClearDescription()
			ic.state.PrevPhase()
		} else {
			// We have a posting to go back to - clear last value and go back
			ic.state.JournalEntryInput.CurrentPosting().ClearValue()
			ic.state.SetPhase(statemod.InputPostingValue)
		}
	case statemod.InputPostingValue:
		ic.state.JournalEntryInput.CurrentPosting().ClearAccount()
		ic.state.PrevPhase()
	default:
	}
}
