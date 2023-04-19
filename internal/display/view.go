package display

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/state"
)

type (
	View struct {
		textView *tview.TextView
		state    *state.State
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

	view := &View{textView: textView, state: state}

	state.AddOnChangeHook(view.refresh)

	return view
}

func (v *View) GetContent() tview.Primitive {
	return v.textView
}

func (v *View) refresh() {
	text := v.state.JournalEntryInput.Repr()
	v.textView.SetText(text)
}
