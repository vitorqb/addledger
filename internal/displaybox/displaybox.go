package displaybox

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/state"
)

type (
	DisplayBox struct {
		textView *tview.TextView
		state    *state.State
	}
)

const (
	BackgroundColor = tcell.ColorBlueViolet
)

func NewDisplayBox(state *state.State) *DisplayBox {
	textView := tview.NewTextView()
	textView.SetBackgroundColor(BackgroundColor)
	textView.SetBorderPadding(1, 1, 1, 1)
	textView.SetBorder(true)
	displayBox := &DisplayBox{textView: textView, state: state}
	state.AddOnChangeHook(displayBox.Refresh)
	return displayBox
}

func (d *DisplayBox) GetTextView() tview.Primitive {
	return d.textView
}

func (d *DisplayBox) Refresh() {
	var text string
	if date, found := d.state.JournalEntryInput.GetDate(); found {
		text += date.Format("2006-01-02")
	}
	if description, found := d.state.JournalEntryInput.GetDescription(); found {
		text += " " + description
	}
	i := -1
	for {
		i++
		if posting, found := d.state.JournalEntryInput.GetPosting(i); found {
			text += "\n" + "    "
			if account, found := posting.GetAccount(); found {
				text += account
			}
		} else {
			break
		}
	}
	d.textView.SetText(text)
}
