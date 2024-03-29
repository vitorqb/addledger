package display

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/userinput"
)

type (
	View struct {
		*tview.TextView
		state *state.State
	}
)

const (
	BackgroundColor = tcell.ColorBlueViolet
)

func NewView(state *state.State) *View {
	textView := tview.NewTextView()
	textView.SetBackgroundColor(BackgroundColor)
	textView.SetBorderPadding(1, 1, 1, 1)
	textView.SetBorder(true)

	view := &View{TextView: textView, state: state}

	state.AddOnChangeHook(view.refresh)

	return view
}

func (v *View) refresh() {
	text := userinput.TransactionRepr(v.state.Transaction)
	v.SetText(text)
}
