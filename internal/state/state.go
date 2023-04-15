package state

import (
	"github.com/vitorqb/addledger/internal/input"
)

type (
	Phase        string
	OnChangeHook func()
	State        struct {
		onChangeHooks     []OnChangeHook
		currentPhase      Phase
		JournalEntryInput *input.JournalEntryInput
	}
)

const (
	InputDate           Phase = "INPUT_DATE"
	InputDescription    Phase = "INPUT_DESCRIPTION"
	InputPostingAccount Phase = "INPUT_POSTING_ACCOUNT"
	InputPostingValue   Phase = "INPUT_POSTING_VALUE"
)

func InitialState() *State {
	journalEntryInput := input.NewJournalEntryInput()
	state := &State{[]OnChangeHook{}, InputDate, journalEntryInput}
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

func (s *State) NextPhase() {
	switch s.currentPhase {
	case InputDate:
		s.currentPhase = InputDescription
	case InputDescription:
		s.currentPhase = InputPostingAccount
	case InputPostingAccount:
		s.currentPhase = InputPostingValue
	default:
	}
	s.notifyChange()
}
