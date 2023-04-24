package state

import (
	"github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/pkg/hledger"
)

type (
	Phase        string
	OnChangeHook func()
	State        struct {
		onChangeHooks     []OnChangeHook
		currentPhase      Phase
		JournalEntryInput *input.JournalEntryInput
		// accounts is an array w/ all known accounts.
		accounts []string
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
	state := &State{[]OnChangeHook{}, InputDate, journalEntryInput, []string{}}
	journalEntryInput.AddOnChangeHook(state.notifyChange)
	return state
}

func (s *State) CurrentPhase() Phase {
	return s.currentPhase
}

func (s *State) SetPhase(p Phase) {
	s.currentPhase = p
	s.notifyChange()
}

func (s *State) AddOnChangeHook(h OnChangeHook) {
	s.onChangeHooks = append(s.onChangeHooks, h)
}

func (s *State) notifyChange() {
	for _, h := range s.onChangeHooks {
		h()
	}
}

func (s *State) SetAccounts(x []string) {
	s.accounts = x
	s.notifyChange()
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
	s.notifyChange()
}

// LoadMetadata loads metadata to state from Hledger
func (s *State) LoadMetadata(hledgerClient hledger.IClient) error {
	accounts, err := hledgerClient.Accounts()
	if err != nil {
		return err
	}
	s.SetAccounts(accounts)
	return nil
}
