package state

import (
	"time"

	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/utils"
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
		selectedPostingAccount string
		descriptionText        string
		selectedDescription    string

		// Controls posting account
		postingAccountText  string
		postingAccountGuess *MaybeValue[journal.Account]

		// Controls posting ammount
		postingAmmountGuess *MaybeValue[finance.Ammount]
		postingAmmountInput *MaybeValue[finance.Ammount]
		postingAmmountText  string

		// Controls tags
		tagsText    string
		selectedTag journal.Tag

		// Controls date
		dateGuess *MaybeValue[time.Time]
		dateText  string

		// The transactions that match the current input
		matchingTransactions []journal.Transaction
	}

	// Display is the state relative to the display.
	Display struct {
		react.IReact

		// Controls whether the shortcut modal is displayed or not
		shortcutModal      bool
		loadStatementModal bool

		// A message to display to the user
		userMessage string
	}

	// State is the top-level app state
	State struct {
		react.IReact
		currentPhase    Phase
		Transaction     *TransactionData
		InputMetadata   *InputMetadata
		JournalMetadata *JournalMetadata
		// StatementEntries are entires loaded from a bank statement.
		// They are used to help the user to create journal entries.
		StatementEntries []finance.StatementEntry
		Display          *Display
	}

	// MaybeValue is a helper container that may contain a value or not
	MaybeValue[T any] struct {
		react.React
		isSet bool
		value T
	}

	// ArrayValue is a helper container that contains an array of values
	ArrayValue[T any] struct {
		react.React
		value []T
	}
)

func (mv *MaybeValue[T]) Get() (T, bool) {
	return mv.value, mv.isSet
}

func (mv *MaybeValue[T]) Set(x T) {
	if r, ok := any(x).(react.IReact); ok {
		r.AddOnChangeHook(mv.NotifyChange)
	}
	mv.value = x
	mv.isSet = true
	mv.NotifyChange()
}

func (mv *MaybeValue[T]) Clear() {
	if mv.isSet {
		var zero T
		mv.value = zero
		mv.isSet = false
		mv.NotifyChange()
	}
}

func (av *ArrayValue[T]) Get() []T {
	if av.value == nil {
		return []T{}
	}
	return av.value
}

func (av *ArrayValue[T]) Set(x []T) {
	for _, v := range x {
		if r, ok := any(v).(react.IReact); ok {
			r.AddOnChangeHook(av.NotifyChange)
		}
	}
	av.value = x
	av.NotifyChange()
}

func (av *ArrayValue[T]) Append(x T) {
	if r, ok := any(x).(react.IReact); ok {
		r.AddOnChangeHook(av.NotifyChange)
	}
	av.value = append(av.value, x)
	av.NotifyChange()
}

func (av *ArrayValue[T]) Clear() {
	if len(av.value) > 0 {
		av.value = []T{}
		av.NotifyChange()
	}
}

func (av *ArrayValue[T]) Pop() {
	if len(av.value) > 0 {
		av.value = av.value[:len(av.value)-1]
		av.NotifyChange()
	}
}

func (av *ArrayValue[T]) Last() (T, bool) {
	if len(av.value) > 0 {
		return av.value[len(av.value)-1], true
	}
	var zero T
	return zero, false
}

const (
	InputDate           Phase = "INPUT_DATE"
	InputDescription    Phase = "INPUT_DESCRIPTION"
	InputTags           Phase = "INPUT_TAGS"
	InputPostingAccount Phase = "INPUT_POSTING_ACCOUNT"
	InputPostingAmmount Phase = "INPUT_POSTING_AMMOUNT"
	Confirmation        Phase = "CONFIRMATION"
)

func InitialState() *State {
	inputMetadata := &InputMetadata{
		IReact:                 react.New(),
		selectedPostingAccount: "",
		descriptionText:        "",
		selectedDescription:    "",
		postingAccountText:     "",
		postingAccountGuess:    &MaybeValue[journal.Account]{},
		postingAmmountGuess:    &MaybeValue[finance.Ammount]{},
		postingAmmountInput:    &MaybeValue[finance.Ammount]{},
		postingAmmountText:     "",
		dateGuess:              &MaybeValue[time.Time]{},
		dateText:               "",
		matchingTransactions:   []journal.Transaction{},
	}
	journalMetadata := NewJournalMetadata()
	display := NewDisplay()
	state := &State{
		IReact:           react.New(),
		currentPhase:     InputDate,
		Transaction:      NewTransactionData(),
		InputMetadata:    inputMetadata,
		JournalMetadata:  journalMetadata,
		StatementEntries: []finance.StatementEntry{},
		Display:          display,
	}
	inputMetadata.AddOnChangeHook(state.NotifyChange)
	journalMetadata.AddOnChangeHook(state.NotifyChange)
	display.AddOnChangeHook(state.NotifyChange)
	state.Transaction.AddOnChangeHook(state.NotifyChange)
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
		s.currentPhase = InputTags
	case InputTags:
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
	case InputTags:
		s.currentPhase = InputDescription
	case InputPostingAccount:
		s.currentPhase = InputTags
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

// GetPostingAccountGuess returns the current guess for the PostingAccount input.
func (im *InputMetadata) GetPostingAccountGuess() (journal.Account, bool) {
	return im.postingAccountGuess.Get()
}

// SetPostingAccountGuess sets the current guess for the PostingAccount input.
func (im *InputMetadata) SetPostingAccountGuess(x journal.Account) {
	im.postingAccountGuess.Set(x)
	im.NotifyChange()
}

// ClearPostingAccountGuess clears the current guess for the PostingAccount input.
func (im *InputMetadata) ClearPostingAccountGuess() {
	im.postingAccountGuess.Clear()
	im.NotifyChange()
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
func (im *InputMetadata) SetPostingAmmountGuess(x finance.Ammount) {
	im.postingAmmountGuess.Set(x)
	im.NotifyChange()
}

// GetPostingAmmountGuess returns the current guess for the ammount to enter. The
// second returned value described whether the value is set or not.
func (im *InputMetadata) GetPostingAmmountGuess() (finance.Ammount, bool) {
	return im.postingAmmountGuess.Get()
}

// ClearPostingAmmountGuess cleats the guess for the ammount to enter.
func (im *InputMetadata) ClearPostingAmmountGuess() {
	im.postingAmmountGuess.Clear()
	im.NotifyChange()
}

// SetPostingAmmountInput sets the current inputted ammount by the user.
func (im *InputMetadata) SetPostingAmmountInput(x finance.Ammount) {
	im.postingAmmountInput.Set(x)
	im.NotifyChange()
}

// GetPostingAmmountInput returns the current input for the ammount to enter. The
// second returned value described whether the value is set or not.
func (im *InputMetadata) GetPostingAmmountInput() (finance.Ammount, bool) {
	return im.postingAmmountInput.Get()
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
	return im.dateGuess.Get()
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

// GetDateText returns the current text for the date input
func (im *InputMetadata) GetDateText() string {
	return im.dateText
}

// SetDateText sets the current text for the date input
func (im *InputMetadata) SetDateText(x string) {
	im.dateText = x
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

// MatchingTransactions returns the current matching transactions
func (im *InputMetadata) MatchingTransactions() []journal.Transaction {
	return im.matchingTransactions
}

// SetMatchingTransactions sets the current matching transactions
func (im *InputMetadata) SetMatchingTransactions(x []journal.Transaction) {
	im.matchingTransactions = x
	im.NotifyChange()
}

// TagText returns the current text for the Tags input
func (im *InputMetadata) TagText() string { return im.tagsText }

// SetTagText sets the current text for the Tags input
func (im *InputMetadata) SetTagText(x string) {
	im.tagsText = x
	im.NotifyChange()
}

// SetSelectedTag sets the current text for the selected tag in the
// context's TagList.
func (im *InputMetadata) SetSelectedTag(x journal.Tag) {
	if im.selectedTag != x {
		im.selectedTag = x
		im.NotifyChange()
	}
}

// SelectedTag returns the current text for the selected tag in the
// context's TagList.
func (im *InputMetadata) SelectedTag() journal.Tag {
	return im.selectedTag
}

// Reset resets all user input from the InputMetadata
func (im *InputMetadata) Reset() {
	im.postingAccountText = ""
	im.selectedPostingAccount = ""
	im.descriptionText = ""
	im.selectedDescription = ""
	im.postingAmmountGuess = &MaybeValue[finance.Ammount]{}
	im.postingAmmountInput = &MaybeValue[finance.Ammount]{}
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

// Tags returns all known tags for the journal
func (jm *JournalMetadata) Tags() []journal.Tag {
	tags := []journal.Tag{}
	for _, transaction := range jm.transactions {
		tags = append(tags, transaction.Tags...)
	}
	return utils.Unique(tags)
}

// GetStatementEntries returns the current statement entries
func (s *State) GetStatementEntries() []finance.StatementEntry {
	return s.StatementEntries
}

// SetStatementEntries sets the current statement entries
func (s *State) SetStatementEntries(x []finance.StatementEntry) {
	s.StatementEntries = x
	s.NotifyChange()
}

// CurrentStatementEntry returns the current statement entry
func (s *State) CurrentStatementEntry() (e finance.StatementEntry, found bool) {
	if len(s.StatementEntries) == 0 {
		return finance.StatementEntry{}, false
	}
	return s.StatementEntries[0], true
}

// PopStatementEntry pops the current statement entry
func (s *State) PopStatementEntry() {
	if len(s.StatementEntries) > 0 {
		s.StatementEntries = s.StatementEntries[1:]
		s.NotifyChange()
	}
}

// NewDisplay returns a new Display
func NewDisplay() *Display {
	return &Display{
		IReact:             react.New(),
		shortcutModal:      false,
		loadStatementModal: false,
		userMessage:        "",
	}
}

// SetShortcutModal sets whether the statement modal is displayed or not
func (d *Display) SetShortcutModal(x bool) {
	d.shortcutModal = x
	d.NotifyChange()
}

// ShortcutModal returns whether the statement modal is displayed or not
func (d *Display) ShortcutModal() bool {
	return d.shortcutModal
}

// SetLoadStatementModal sets whether the statement modal is displayed or not
func (d *Display) SetLoadStatementModal(x bool) {
	d.loadStatementModal = x
	d.NotifyChange()
}

// LoadStatementModal returns whether the statement modal is displayed or not
func (d *Display) LoadStatementModal() bool {
	return d.loadStatementModal
}

// SetUserMessage sets the current user message
func (d *Display) SetUserMessage(x string) {
	d.userMessage = x
	d.NotifyChange()
}

// UserMessage returns the current user message
func (d *Display) UserMessage() string {
	return d.userMessage
}
