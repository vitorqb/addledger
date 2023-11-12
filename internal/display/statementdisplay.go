package display

import (
	"fmt"

	"github.com/rivo/tview"
	statemod "github.com/vitorqb/addledger/internal/state"
)

type StatementDisplay struct {
	*tview.TextView
	state *statemod.State
}

func NewStatementDisplay(state *statemod.State) *StatementDisplay {
	out := &StatementDisplay{TextView: tview.NewTextView(), state: state}
	out.SetBorder(true)

	out.Refresh()
	state.AddOnChangeHook(out.Refresh)
	return out
}

func (s *StatementDisplay) Refresh() {
	staEntry, found := s.state.CurrentStatementEntry()
	if !found {
		s.SetText("")
		return
	}
	s.SetText(fmt.Sprintf(
		"%s | %s | %s | %s %s | [%d]",
		staEntry.Date.Format("2006/01/02"),
		staEntry.Description,
		staEntry.Account,
		staEntry.Ammount.Commodity,
		staEntry.Ammount.Quantity.String(),
		len(s.state.StatementEntries),
	))
}
