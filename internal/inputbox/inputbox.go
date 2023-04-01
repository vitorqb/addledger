package inputbox

import (
	"github.com/rivo/tview"
)

func NewInputBox() *tview.Box {
	return tview.NewBox().SetBorder(true).SetTitle("Input Area")
}
