package displaybox

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/utils"
)

type (
	DisplayBox struct {
		textView *tview.TextView
		date      time.Time
		description string
	}
)

func NewDisplayBox() *DisplayBox {
	textView := tview.NewTextView()
	textView.SetBackgroundColor(tcell.ColorBlueViolet)
	textView.SetBorderPadding(1, 1, 1, 1)
	return &DisplayBox{textView: textView}
}

func (d *DisplayBox) GetTextView() tview.Primitive {
	return utils.Center(100, 15, d.textView)
}

func (d *DisplayBox) SetDate(x time.Time) {
	d.date = x
	d.refresh()
}

func (d *DisplayBox) SetDescription(x string) {
	d.description = x
	d.refresh()
}

func (d *DisplayBox) refresh() {
	date := d.date.Format("2006-01-02")
	description := d.description
	d.textView.SetText(date + " " + description)
}
