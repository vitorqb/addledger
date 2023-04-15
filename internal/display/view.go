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
	// !!! TODO Make a Repr() method on JournalEntryInput
	var text string
	if date, found := v.state.JournalEntryInput.GetDate(); found {
		text += date.Format("2006-01-02")
	}
	if description, found := v.state.JournalEntryInput.GetDescription(); found {
		text += " " + description
	}
	i := -1
	for {
		i++
		if posting, found := v.state.JournalEntryInput.GetPosting(i); found {
			text += "\n" + "    "
			if account, found := posting.GetAccount(); found {
				text += account
			}
			text += "    "
			if value, found := posting.GetValue(); found {
				text += value
			}
		} else {
			break
		}
	}
	v.textView.SetText(text)
}
