package context

import (
	"github.com/rivo/tview"
	statemod "github.com/vitorqb/addledger/internal/state"
)

type AmmountGuesser struct {
	*tview.TextView
	state *statemod.State
}

func NewAmmountGuesser(state *statemod.State) (*AmmountGuesser, error) {
	guesser := &AmmountGuesser{tview.NewTextView(), state}
	guesser.Refresh()
	state.AddOnChangeHook(guesser.Refresh)
	return guesser, nil
}

func (ag *AmmountGuesser) Refresh() {
	guess, found := ag.state.InputMetadata.GetPostingAmmountGuess()
	if found {
		newText := guess.Commodity + " " + guess.Quantity.String()
		ag.SetText(newText)
	} else {
		ag.SetText("")
	}
}
