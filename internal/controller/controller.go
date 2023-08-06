package controller

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/dateguesser"
	"github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/listaction"
	statemod "github.com/vitorqb/addledger/internal/state"
)

//go:generate $MOCKGEN --source=controller.go --destination=../../mocks/controller/controller_mock.go

// IInputController reacts to the user inputs and interactions.
type IInputController interface {
	// Handles user confirming/rejecting the transaction at the end, after
	// inputing everything
	OnInputConfirmation()
	OnInputRejection()

	// Handles user entering a new date for a posting
	OnDateChanged(text string)
	OnDateDone()

	// Handles user entering a new ammount for a posting
	OnPostingAmmountChanged(text string)
	OnPostingAmmountDone(input.DoneSource)

	// Called when an user wants to undo it's last action.
	OnUndo()

	// Controls Posting Account input
	OnPostingAccountChanged(newText string)
	OnPostingAccountDone(source input.DoneSource)
	OnPostingAccountInsertFromContext()
	OnPostingAccountListAcction(action listaction.ListAction)

	// Controls the Description input
	OnDescriptionChanged(newText string)
	OnDescriptionDone()
	OnDescriptionInsertFromContext()
	OnDescriptionListAction(action listaction.ListAction)
	OnDescriptionSelectedFromContext()
}

// InputController implements IInputController.
type InputController struct {
	state       *statemod.State
	output      io.Writer
	eventBus    eventbus.IEventBus
	dateGuesser dateguesser.IDateGuesser
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
	if opts.dateGuesser == nil {
		return nil, fmt.Errorf("missing DateGuesser")
	}
	return &InputController{
		state:       state,
		output:      opts.output,
		eventBus:    opts.eventBus,
		dateGuesser: opts.dateGuesser,
	}, nil
}

func (ic *InputController) OnDateChanged(x string) {
	date, success := ic.dateGuesser.Guess(x)
	if success {
		ic.state.InputMetadata.SetDateGuess(date)
	} else {
		ic.state.InputMetadata.ClearDateGuess()
	}
}

func (ic *InputController) OnDateDone() {
	if date, found := ic.state.InputMetadata.GetDateGuess(); found {
		ic.state.JournalEntryInput.SetDate(date)
		ic.state.NextPhase()
	}
}

func (ic *InputController) OnPostingAccountDone(source input.DoneSource) {
	account := ""
	switch source {
	case input.Context:
		account = ic.state.InputMetadata.SelectedPostingAccount()
	case input.Input:
		account = ic.state.InputMetadata.PostingAccountText()
	}

	// If account is empty - do nothing
	if account == "" {
		return
	}

	// We have an account - save the posting
	posting := ic.state.JournalEntryInput.CurrentPosting()
	posting.SetAccount(account)

	// Go to ammount
	ic.state.NextPhase()
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

func (ic *InputController) OnPostingAmmountDone(source input.DoneSource) {
	var ammount journal.Ammount
	var success bool

	switch source {
	case input.Context:
		ammount, success = ic.state.InputMetadata.GetPostingAmmountGuess()
	case input.Input:
		ammount, success = ic.state.InputMetadata.GetPostingAmmountInput()
	default:
		logrus.Fatalf("Uknown source for Posting Ammount: %s", source)
	}

	if success {
		// Saves ammount
		posting := ic.state.JournalEntryInput.CurrentPosting()
		posting.SetAmmount(ammount)

		// If there is balance outstanding, go to next posting
		if !ic.state.JournalEntryInput.PostingHasZeroBalance() {
			ic.state.JournalEntryInput.AdvancePosting()
			ic.state.SetPhase(statemod.InputPostingAccount)
			return
		}

		// Else, go to confirmation
		ic.state.SetPhase(statemod.Confirmation)
	}
}

func (ic *InputController) OnPostingAmmountChanged(text string) {
	if text != ic.state.InputMetadata.GetPostingAmmountText() {
		ic.state.InputMetadata.SetPostingAmmountText(text)
		ammount, err := input.TextToAmmount(text)
		if err != nil {
			ic.state.InputMetadata.ClearPostingAmmountInput()
		} else {
			ic.state.InputMetadata.SetPostingAmmountInput(ammount)
		}
	}
}

func (ic *InputController) OnInputConfirmation() {
	_, err := io.WriteString(ic.output, "\n\n"+ic.state.JournalEntryInput.Repr())
	if err != nil {
		// TODO Let user know somehow!
		logrus.WithError(err).Fatal("failed to write to file")
		return
	}
	ic.state.JournalEntryInput = input.NewJournalEntryInput()
	ic.state.InputMetadata.Reset()
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
			ic.state.InputMetadata.SetDescriptionText("")
			ic.state.PrevPhase()
		} else {
			// We have a posting to go back to - clear last ammount and go back
			ic.state.JournalEntryInput.CurrentPosting().ClearAmmount()
			ic.state.SetPhase(statemod.InputPostingAmmount)
		}
	case statemod.InputPostingAmmount:
		ic.state.JournalEntryInput.CurrentPosting().ClearAccount()
		ic.state.PrevPhase()
	default:
	}
}
