package displaybox

import (
	"github.com/rivo/tview"
)

func NewDisplayBox() *tview.Box {
	return tview.NewBox().SetBorder(true).SetTitle("Display Area")
}
