package state

import (
	"github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/pkg/hledger"
	"github.com/vitorqb/addledger/pkg/react"
)

type (
	Phase string

	// JournalMetadata is the state relative to the journal.
	JournalMetadata struct {
		react.IReact
		// transactions is a list with all transactions
		transactions []journal.Transaction
	}

	// InputMetadata is the state relative to inputs.
	InputMetadata struct {
		react.IReact
		postingAccountText     string
		selectedPostingAccount string
		descriptionText        string
		selectedDescription    string
	}

	// State is the top-level app state
	State struct {
		react.IReact
		currentPhase      Phase
		JournalEntryInput *input.JournalEntryInput
		// accounts is an array w/ all known accounts.
		// !!!! TODO Move accounts to JournalMetadata
		accounts        []string
		InputMetadata   *InputMetadata
		JournalMetadata *JournalMetadata
	}
)

const (
	InputDate           Phase = "INPUT_DATE"
	InputDescription    Phase = "INPUT_DESCRIPTION"
	InputPostingAccount Phase = "INPUT_POSTING_ACCOUNT"
	InputPostingValue   Phase = "INPUT_POSTING_VALUE"
	Confirmation        Phase = "CONFIRMATION"
)

func InitialState() *State {
	journalEntryInput := input.NewJournalEntryInput()
	inputMetadata := &InputMetadata{react.New(), "", "", "", ""}
	// !!!! TODO Add postings here, or have a setup method that loads them.
	journalMetadata := &JournalMetadata{react.New(), []journal.Transaction{}}
	state := &State{
		react.New(),
		InputDate,
		journalEntryInput,
		[]string{},
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

func (s *State) SetAccounts(x []string) {
	s.accounts = x
	s.NotifyChange()
}

func (s *State) GetAccounts() []string {
	return s.accounts
}

func (s *State) NextPhase() {
	switch s.currentPhase {
	case InputDate:
		s.currentPhase = InputDescription
	case InputDescription:
		s.currentPhase = InputPostingAccount
	case InputPostingAccount:
		s.currentPhase = InputPostingValue
	case InputPostingValue:
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
	case InputPostingValue:
		s.currentPhase = InputPostingAccount
	case Confirmation:
		s.currentPhase = InputPostingValue
	default:
	}
	s.NotifyChange()
}

// LoadMetadata loads metadata to state from Hledger
func (s *State) LoadMetadata(hledgerClient hledger.IClient) error {

	// load accounts
	accounts, err := hledgerClient.Accounts()
	if err != nil {
		return err
	}
	s.SetAccounts(accounts)

	// load postings
	postings, err := hledgerClient.Transactions()
	if err != nil {
		return err
	}
	s.JournalMetadata.SetTransactions(postings)

	return nil
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

// DescriptionText returns the current text for the selected description in
// the context.
func (im *InputMetadata) DescriptionText() string { return im.descriptionText }

// SetDescriptionText sets the current text for the selected description in
// the context.
func (im *InputMetadata) SetDescriptionText(x string) {
	im.descriptionText = x
	im.NotifyChange()
}

// Transactions returns all known postings for the journal
func (jm *JournalMetadata) Transactions() []journal.Transaction { return jm.transactions }

// SetTransactions sets all known postings for the journal
func (jm *JournalMetadata) SetTransactions(x []journal.Transaction) {
	jm.transactions = x
	jm.NotifyChange()
}
