package state

import (
	"time"

	"github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/pkg/react"
)

type (
	Phase string

	// JournalMetadata is the state relative to the journal.
	JournalMetadata struct {
		react.IReact
		// transactions is a list with all transactions
		transactions []journal.Transaction
		// accounts is a list of all known accounts
		accounts []journal.Account
	}

	// InputMetadata is the state relative to inputs.
	InputMetadata struct {
		react.IReact
		postingAccountText     string
		selectedPostingAccount string
		descriptionText        string
		selectedDescription    string

		// Controls posting ammount
		postingAmmountGuess *MaybeValue[journal.Ammount]
		postingAmmountInput *MaybeValue[journal.Ammount]
		postingAmmountText  string

		// Controls date
		dateGuess *MaybeValue[time.Time]
	}

	// State is the top-level app state
	State struct {
		react.IReact
		currentPhase      Phase
		JournalEntryInput *input.JournalEntryInput
		InputMetadata     *InputMetadata
		JournalMetadata   *JournalMetadata
	}

	// MaybeValue is a helper container that may contain a value or not
	MaybeValue[T any] struct {
		IsSet bool
		Value T
	}
)

func (mv *MaybeValue[T]) Set(x T) {
	mv.Value = x
	mv.IsSet = true
}

func (mv *MaybeValue[T]) Clear() {
	mv.IsSet = false
}

const (
	InputDate           Phase = "INPUT_DATE"
	InputDescription    Phase = "INPUT_DESCRIPTION"
	InputPostingAccount Phase = "INPUT_POSTING_ACCOUNT"
	InputPostingAmmount Phase = "INPUT_POSTING_AMMOUNT"
	Confirmation        Phase = "CONFIRMATION"
)

func InitialState() *State {
	journalEntryInput := input.NewJournalEntryInput()
	inputMetadata := &InputMetadata{
		IReact:                 react.New(),
		postingAccountText:     "",
		selectedPostingAccount: "",
		descriptionText:        "",
		selectedDescription:    "",
		postingAmmountGuess:    &MaybeValue[journal.Ammount]{},
		postingAmmountInput:    &MaybeValue[journal.Ammount]{},
		postingAmmountText:     "",
		dateGuess:              &MaybeValue[time.Time]{},
	}
	journalMetadata := NewJournalMetadata()
	state := &State{
		react.New(),
		InputDate,
		journalEntryInput,
		inputMetadata,
		journalMetadata,
	}
	journalEntryInput.AddOnChangeHook(state.NotifyChange)
	inputMetadata.AddOnChangeHook(state.NotifyChange)
	journalMetadata.AddOnChangeHook(state.NotifyChange)
	return state
}

func (s *State) CurrentPhase() Phase {
	return s.currentPhase
}

func (s *State) SetPhase(p Phase) {
	s.currentPhase = p
	s.NotifyChange()
}

func (s *State) NextPhase() {
	switch s.currentPhase {
	case InputDate:
		s.currentPhase = InputDescription
	case InputDescription:
		s.currentPhase = InputPostingAccount
	case InputPostingAccount:
		s.currentPhase = InputPostingAmmount
	case InputPostingAmmount:
		s.currentPhase = Confirmation
	default:
	}
	s.NotifyChange()
}

func (s *State) PrevPhase() {
	switch s.currentPhase {
	case InputDescription:
		s.currentPhase = InputDate
	case InputPostingAccount:
		s.currentPhase = InputDescription
	case InputPostingAmmount:
		s.currentPhase = InputPostingAccount
	case Confirmation:
		s.currentPhase = InputPostingAmmount
	default:
	}
	s.NotifyChange()
}

// PostingAccountText returns the current text for the PostingAccount input.
func (im *InputMetadata) PostingAccountText() string {
	return im.postingAccountText
}

// PostingAccountText sets the current text for the PostingAccount input.
func (im *InputMetadata) SetPostingAccountText(x string) {
	if im.postingAccountText != x {
		im.postingAccountText = x
		im.NotifyChange()
	}
}

// SelectedPostingAccount returns the current text for the selected account in the
// context's AccountList.
func (im *InputMetadata) SelectedPostingAccount() string {
	return im.selectedPostingAccount
}

// SetSelectedPostingAccount sets the current text for the selected account in the
// context's AccountList.
func (im *InputMetadata) SetSelectedPostingAccount(x string) {
	if im.selectedPostingAccount != x {
		im.selectedPostingAccount = x
		im.NotifyChange()
	}
}

// SelectedDescription return the current text for the selected description in
// the context.
func (im *InputMetadata) SelectedDescription() string {
	return im.selectedDescription
}

// SetSelectedDescription sets the current text for the selected description in
// the context.
func (im *InputMetadata) SetSelectedDescription(x string) {
	if im.selectedDescription != x {
		im.selectedDescription = x
		im.NotifyChange()
	}
}

// SetPostingAmmountGuess sets the current guess for the ammount to enter.
func (im *InputMetadata) SetPostingAmmountGuess(x journal.Ammount) {
	im.postingAmmountGuess.Set(x)
	im.NotifyChange()
}

// GetPostingAmmountGuess returns the current guess for the ammount to enter. The
// second returned value described whether the value is set or not.
func (im *InputMetadata) GetPostingAmmountGuess() (journal.Ammount, bool) {
	if !im.postingAmmountGuess.IsSet {
		return journal.Ammount{}, false
	}
	return im.postingAmmountGuess.Value, true
}

// ClearPostingAmmountGuess cleats the guess for the ammount to enter.
func (im *InputMetadata) ClearPostingAmmountGuess() {
	im.postingAmmountGuess.Clear()
	im.NotifyChange()
}

// SetPostingAmmountInput sets the current inputted ammount by the user.
func (im *InputMetadata) SetPostingAmmountInput(x journal.Ammount) {
	im.postingAmmountInput.Set(x)
	im.NotifyChange()
}

// GetPostingAmmountInput returns the current input for the ammount to enter. The
// second returned value described whether the value is set or not.
func (im *InputMetadata) GetPostingAmmountInput() (journal.Ammount, bool) {
	if !im.postingAmmountInput.IsSet {
		return journal.Ammount{}, false
	}
	return im.postingAmmountInput.Value, true
}

// ClearPostingAmmountInput cleats the input for the ammount to enter.
func (im *InputMetadata) ClearPostingAmmountInput() {
	im.postingAmmountInput.Clear()
	im.NotifyChange()
}

// GetPostingAmmountText returns the current text inputted by the user for PostingAmmount.
func (im *InputMetadata) GetPostingAmmountText() string {
	return im.postingAmmountText
}

// SetPostingAmmountText sets the current text inputted by the user for PostingAmmount.
func (im *InputMetadata) SetPostingAmmountText(x string) {
	im.postingAmmountText = x
	im.NotifyChange()
}

// ClearPostingAmmountText sets the current text inputted by the user for PostingAmmount.
func (im *InputMetadata) ClearPostingAmmountText() {
	im.postingAmmountText = ""
	im.NotifyChange()
}

// GetDateGuess returns the current date guess
func (im *InputMetadata) GetDateGuess() (time.Time, bool) {
	if !im.dateGuess.IsSet {
		return time.Time{}, false
	}
	return im.dateGuess.Value, true
}

// SetDateGuess sets the current date guess
func (im *InputMetadata) SetDateGuess(x time.Time) {
	im.dateGuess.Set(x)
	im.NotifyChange()
}

// ClearDateGuess clears the current date guess
func (im *InputMetadata) ClearDateGuess() {
	im.dateGuess.Clear()
	im.NotifyChange()
}

// DescriptionText returns the current text for the selected description in
// the context.
func (im *InputMetadata) DescriptionText() string { return im.descriptionText }

// SetDescriptionText sets the current text for the selected description in
// the context.
func (im *InputMetadata) SetDescriptionText(x string) {
	im.descriptionText = x
	im.NotifyChange()
}

// Reset resets all user input from the InputMetadata
func (im *InputMetadata) Reset() {
	im.postingAccountText = ""
	im.selectedPostingAccount = ""
	im.descriptionText = ""
	im.selectedDescription = ""
	im.postingAmmountGuess = &MaybeValue[journal.Ammount]{}
	im.postingAmmountInput = &MaybeValue[journal.Ammount]{}
	im.postingAmmountText = ""
	im.dateGuess = &MaybeValue[time.Time]{}
	im.NotifyChange()
}

func NewJournalMetadata() *JournalMetadata {
	return &JournalMetadata{
		react.New(),
		[]journal.Transaction{},
		[]journal.Account{},
	}
}

// Transactions returns all known postings for the journal
func (jm *JournalMetadata) Transactions() []journal.Transaction { return jm.transactions }

// SetTransactions sets all known postings for the journal
func (jm *JournalMetadata) SetTransactions(x []journal.Transaction) {
	jm.transactions = x
	jm.NotifyChange()
}

// AppendTransaction appends a transaction to the JournalMetadata.
func (jm *JournalMetadata) AppendTransaction(x journal.Transaction) {
	jm.transactions = append(jm.transactions, x)
	jm.NotifyChange()
}

// Accounts returns all known postings for the journal
func (jm *JournalMetadata) Accounts() []journal.Account { return jm.accounts }

// SetAccounts sets all known postings for the journal
func (jm *JournalMetadata) SetAccounts(x []journal.Account) {
	jm.accounts = x
	jm.NotifyChange()
}
