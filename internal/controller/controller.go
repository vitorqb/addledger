package controller

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/dateguesser"
	"github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/listaction"
	"github.com/vitorqb/addledger/internal/metaloader"
	printermod "github.com/vitorqb/addledger/internal/printer"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/userinput"
	"github.com/vitorqb/addledger/internal/usermessenger"
)

//go:generate $MOCKGEN --source=controller.go --destination=../../mocks/controller/controller_mock.go

// StatementLoader represents a component that loads a statement into the app state.
type StatementLoader interface {
	LoadFromFiles(statementFile, presetFile string) error
}

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
	OnPostingAmmountDone(userinput.DoneSource)

	// Called when an user wants to undo it's last action.
	OnUndo()

	// Controls Posting Account input
	OnPostingAccountChanged(newText string)
	OnPostingAccountDone(source userinput.DoneSource)
	OnPostingAccountInsertFromContext()
	OnPostingAccountListAcction(action listaction.ListAction)
	OnFinishPosting()

	// Controls the Description input
	OnDescriptionChanged(newText string)
	OnDescriptionDone(source userinput.DoneSource)
	OnDescriptionInsertFromContext()
	OnDescriptionListAction(action listaction.ListAction)

	// Controls the Tags input
	OnTagChanged(newText string)
	OnTagDone(source userinput.DoneSource)
	OnTagInsertFromContext()
	OnTagListAction(action listaction.ListAction)

	// Controls statement
	OnLoadStatement(csvFile string, presetFile string)
	OnDiscardStatement()
	OnLoadStatementRequest()

	// Controls shortcuts modal
	OnDisplayShortcutModal()
	OnHideShortcutModal()
	OnShowStatementModal()
	OnHideStatementModal()
}

// InputController implements IInputController.
type InputController struct {
	state              *statemod.State
	output             io.Writer
	eventBus           eventbus.IEventBus
	dateGuesser        dateguesser.IDateGuesser
	metaLoader         metaloader.IMetaLoader
	printer            printermod.IPrinter
	csvStatementLoader StatementLoader
	userMessenger      usermessenger.IUserMessenger
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
	if opts.metaLoader == nil {
		return nil, fmt.Errorf("missing IMetaLoader")
	}
	if opts.printer == nil {
		return nil, fmt.Errorf("missing printer")
	}
	if opts.csvStatementLoader == nil {
		return nil, fmt.Errorf("missing csvStatementLoader")
	}
	if opts.userMessenger == nil {
		opts.userMessenger = &usermessenger.NoOp{}
	}
	return &InputController{
		state:              state,
		output:             opts.output,
		eventBus:           opts.eventBus,
		dateGuesser:        opts.dateGuesser,
		metaLoader:         opts.metaLoader,
		printer:            opts.printer,
		csvStatementLoader: opts.csvStatementLoader,
		userMessenger:      opts.userMessenger,
	}, nil
}

func (ic *InputController) OnDateChanged(x string) {
	ic.state.InputMetadata.SetDateText(x)
}

func (ic *InputController) OnDateDone() {
	if date, found := ic.state.InputMetadata.GetDateGuess(); found {
		ic.state.Transaction.Date.Set(date)
		ic.state.NextPhase()
		return
	}
}

func (ic *InputController) OnPostingAccountDone(source userinput.DoneSource) {
	account := ""
	switch source {
	case userinput.Context:
		account = ic.state.InputMetadata.SelectedPostingAccount()
	case userinput.Input:
		account = ic.state.InputMetadata.PostingAccountText()
	}

	// If account is empty - do nothing
	if account == "" {
		return
	}

	// We have an account - save the posting
	posting, found := ic.state.Transaction.Postings.Last()
	if !found {
		posting = statemod.NewPostingData()
		ic.state.Transaction.Postings.Append(posting)
	}
	posting.Account.Set(journal.Account(account))

	// Go to ammount
	ic.state.NextPhase()
}

// OnFinishPosting may be called by the user to signal it is done entering
// postings. This is useful when the user is entering a transaction with
// multiple commodities, since we can't know when the user is done entering
// based on the pending balance.
func (ic *InputController) OnFinishPosting() {
	balance := userinput.PostingBalance(ic.state.Transaction.Postings.Get())

	// len(balance) == 0 -> No postings, nothing to do
	if len(balance.Ammounts()) == 0 {
		return
	}

	// If we have a single currency and it's unbalanced, nothing to do
	if balance.SingleCommodity() && !balance.IsZero() {
		return
	}

	// We have multiple currencies or a single currency with zero balance, so
	// finish the posting
	ic.state.Transaction.Postings.Pop()
	ic.state.SetPhase(statemod.Confirmation)
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
		ic.userMessenger.Error("Failed to send event", err)
	}
}

func (ic *InputController) OnPostingAccountListAcction(action listaction.ListAction) {
	err := ic.eventBus.Send(eventbus.Event{
		Topic: "input.postingaccount.listaction",
		Data:  action,
	})
	if err != nil {
		ic.userMessenger.Error("Failed to send event", err)
	}
}

func (ic *InputController) OnPostingAccountChanged(newText string) {
	ic.state.InputMetadata.SetPostingAccountText(newText)
}

func (ic *InputController) OnPostingAmmountDone(source userinput.DoneSource) {
	var ammount finance.Ammount
	var success bool

	switch source {
	case userinput.Context:
		ammount, success = ic.state.InputMetadata.GetPostingAmmountGuess()
	case userinput.Input:
		ammount, success = ic.state.InputMetadata.GetPostingAmmountInput()
	default:
		logrus.Fatalf("Uknown source for Posting Ammount: %s", source)
	}

	if success {
		// Saves ammount
		posting, found := ic.state.Transaction.Postings.Last()
		if !found {
			posting = statemod.NewPostingData()
			ic.state.Transaction.Postings.Append(posting)
		}
		posting.Ammount.Set(ammount)

		// If there is balance outstanding, go to next posting
		balance := userinput.PostingBalance(ic.state.Transaction.Postings.Get())
		if !balance.IsZero() {
			newPosting := statemod.NewPostingData()
			ic.state.Transaction.Postings.Append(newPosting)
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
		ammount, err := userinput.TextToAmmount(text)
		if err != nil {
			ic.state.InputMetadata.ClearPostingAmmountInput()
		} else {
			ic.state.InputMetadata.SetPostingAmmountInput(ammount)
		}
	}
}

func (ic *InputController) OnInputConfirmation() {
	transaction, transactionErr := userinput.TransactionFromData(ic.state.Transaction)
	if transactionErr != nil {
		ic.userMessenger.Error("The transaction input could not be parsed (this shouldn't happen)", transactionErr)
		return
	}

	// TODO Inject the printer instead of hardcoding
	printErr := printermod.New(2, 0).Print(ic.output, transaction)
	if printErr != nil {
		ic.userMessenger.Error("Failed to write to file", printErr)
		return
	}
	ic.state.Transaction = statemod.NewTransactionData()
	ic.state.InputMetadata.Reset()
	accountLoadErr := ic.metaLoader.LoadAccounts()
	if accountLoadErr != nil {
		ic.userMessenger.Error("Failed to load accounts", accountLoadErr)
		return
	}
	// Note: we could call `ic.metaLoader.LoadTransactions` here. This is, however,
	// quite slow for large journals.
	ic.state.JournalMetadata.AppendTransaction(transaction)

	// We are done w/ the current loaded statement
	ic.state.PopStatementEntry()

	// Go back to first phase and ensure date is cleared
	ic.state.SetPhase(statemod.InputDate)
	ic.OnDateChanged("")
}

func (ic *InputController) OnInputRejection() {
	// put back an empty posting so the user can add to it
	newPosting := statemod.NewPostingData()
	ic.state.Transaction.Postings.Append(newPosting)
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
		ic.userMessenger.Error("Failed to send event", err)
	}
}

func (ic *InputController) OnDescriptionDone(source userinput.DoneSource) {
	if source == userinput.Context {
		// If we have a description from context, use it!
		if descriptionFromContext := ic.state.InputMetadata.SelectedDescription(); descriptionFromContext != "" {
			ic.OnDescriptionChanged(descriptionFromContext)
		}
	}

	description := ic.state.InputMetadata.DescriptionText()
	ic.state.Transaction.Description.Set(description)
	if len(ic.state.Transaction.Postings.Get()) == 0 {
		newPosting := statemod.NewPostingData()
		ic.state.Transaction.Postings.Append(newPosting)
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
		ic.userMessenger.Error("Failed to send event", err)
	}
}

func (ic *InputController) OnTagChanged(newText string) {
	ic.state.InputMetadata.SetTagText(newText)
}

func (ic *InputController) OnTagDone(source userinput.DoneSource) {
	var tag journal.Tag

	// Get tag value
	if source == userinput.Context {
		tag = ic.state.InputMetadata.SelectedTag()
	}
	if tag.Name == "" {
		tag, _ = userinput.TextToTag(ic.state.InputMetadata.TagText())
	}

	// No tag and empty input - move to next phase
	if tag.Name == "" && ic.state.InputMetadata.TagText() == "" {
		ic.state.NextPhase()
		return
	}

	// Skip if no tag - user entered invalid input
	if tag.Name == "" {
		return
	}

	// We have tag - move on to next tag
	ic.state.Transaction.Tags.Append(tag)
	ic.state.InputMetadata.SetTagText("")
	err := ic.eventBus.Send(eventbus.Event{
		Topic: "input.tag.settext",
		Data:  "",
	})
	if err != nil {
		ic.userMessenger.Error("Failed to send event", err)
	}
}

func (ic *InputController) OnTagInsertFromContext() {
	tagFromContext := ic.state.InputMetadata.SelectedTag()
	textFromContext := userinput.TagToText(tagFromContext)
	event := eventbus.Event{
		Topic: "input.tag.settext",
		Data:  textFromContext,
	}
	err := ic.eventBus.Send(event)
	if err != nil {
		ic.userMessenger.Error("Failed to send event", err)
	}
}

func (ic *InputController) OnTagListAction(action listaction.ListAction) {
	event := eventbus.Event{
		Topic: "input.tag.listaction",
		Data:  action,
	}
	err := ic.eventBus.Send(event)
	if err != nil {
		ic.userMessenger.Error("Failed to send event", err)
	}
}

func (ic *InputController) OnDisplayShortcutModal() {
	ic.state.Display.SetShortcutModal(true)
}

// OnHideShortcutModal implements IInputController.
func (ic *InputController) OnHideShortcutModal() {
	ic.state.Display.SetShortcutModal(false)
}

// OnDiscardStatement implements IInputController.
func (ic *InputController) OnDiscardStatement() {
	ic.state.PopStatementEntry()
}

// OnLoadStatementRequest implements IInputController.
func (ic *InputController) OnLoadStatementRequest() {
	ic.state.Display.SetLoadStatementModal(true)
}

func (ic *InputController) OnShowStatementModal() {
	ic.state.Display.StatementModal.SetVisible(true)
}

func (ic *InputController) OnHideStatementModal() {
	ic.state.Display.StatementModal.SetVisible(false)
}

// OnLoadStatement implements display.LoadStatementModalController.
func (ic *InputController) OnLoadStatement(csvFile string, presetFile string) {
	err := ic.csvStatementLoader.LoadFromFiles(csvFile, presetFile)
	if err != nil {
		ic.userMessenger.Error("Failed to load statement", err)
		return
	}
	ic.state.Display.SetLoadStatementModal(false)
}

func (ic *InputController) OnUndo() {
	switch ic.state.CurrentPhase() {
	case statemod.InputDate:
		ic.state.PrevPhase()
	case statemod.InputDescription:
		ic.state.Transaction.Date.Clear()
		ic.state.PrevPhase()
	case statemod.InputTags:
		// Clear description, tags and go back
		ic.state.Transaction.Description.Clear()
		ic.state.InputMetadata.SetDescriptionText("")
		ic.state.Transaction.Tags.Clear()
		ic.state.PrevPhase()
	case statemod.InputPostingAccount:
		ic.state.Transaction.Postings.Pop()
		if posting, found := ic.state.Transaction.Postings.Last(); found {
			// We have a posting to go back to - clear last ammount and go back
			posting.Ammount.Clear()
			ic.state.SetPhase(statemod.InputPostingAmmount)
		} else {
			// We don't have any postings - clear tags and go back
			ic.state.Transaction.Tags.Clear()
			ic.state.InputMetadata.SetTagText("")
			ic.state.PrevPhase()

		}
	case statemod.InputPostingAmmount:
		if posting, found := ic.state.Transaction.Postings.Last(); found {
			posting.Account.Clear()
		}
		ic.state.PrevPhase()
	default:
	}
}
