package state

type (
	Phase string
	State struct {
		currentPhase Phase
	}
)

const (
	Date        Phase = "DATE"
	Description Phase = "DESCRIPTION"
)

func InitialState() *State {
	return &State{Date}
}
