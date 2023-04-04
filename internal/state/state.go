package state

type (
	Phase string
	State struct {
		CurrentPhase Phase
	}
)

const (
	Date        Phase = "DATE"
	Description Phase = "DESCRIPTION"
)

func InitialState() *State {
	return &State{Date}
}

func (s *State) NextPhase() {
	switch s.CurrentPhase {
	case Date:
		s.CurrentPhase = Description
	case Description:
	default:
	}
}
