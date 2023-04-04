package displaybox

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type (
	DisplayBox struct {
		textView    *tview.TextView
		date        time.Time
		description string
	}
)

const (
	BackgroundColor = tcell.ColorBlueViolet
)

func NewDisplayBox() *DisplayBox {
	textView := tview.NewTextView()
	textView.SetBackgroundColor(BackgroundColor)
	textView.SetBorderPadding(1, 1, 1, 1)
	textView.SetBorder(true)
	return &DisplayBox{textView: textView}
}

func (d *DisplayBox) GetTextView() tview.Primitive {
	return d.textView
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
