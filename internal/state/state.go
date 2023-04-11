package state

import (
	"github.com/vitorqb/addledger/internal/input"
)

type (
	Phase        string
	OnChangeHook func()
	State        struct {
		onChangeHooks     []OnChangeHook
		CurrentPhase      Phase
		JournalEntryInput *input.JournalEntryInput
	}
)

const (
	InputDate        Phase = "INPUT_DATE"
	InputDescription Phase = "INPUT_DESCRIPTION"
	InputPostings    Phase = "INPUT_POSTINGS"
)

func InitialState() *State {
	journalEntryInput := input.NewJournalEntryInput()
	state := &State{[]OnChangeHook{}, InputDate, journalEntryInput}
	journalEntryInput.AddOnChangeHook(state.notifyChange)
	return state
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
	switch s.CurrentPhase {
	case InputDate:
		s.CurrentPhase = InputDescription
	case InputDescription:
		s.CurrentPhase = InputPostings
	default:
	}
	s.notifyChange()
}
