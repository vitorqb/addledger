package displaybox

import (
	"time"

	"github.com/rivo/tview"
)

type (
	DisplayBox struct {
		textView *tview.TextView
		date      time.Time
	}
)

func NewDisplayBox() *DisplayBox {
	return &DisplayBox{
		textView: tview.NewTextView(),
	}
}

func (d *DisplayBox) GetPrimitive() tview.Primitive {
	return d.textView
}

func (d *DisplayBox) SetDate(x time.Time) {
	d.date = x
	d.textView.SetText(x.Format("2006-01-02"))
}
