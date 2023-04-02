package main

import (
	"time"

	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/displaybox"
	"github.com/vitorqb/addledger/internal/inputbox"
)

func main() {
	app := tview.NewApplication()
	displayBox := displaybox.NewDisplayBox()
	inputBox := inputbox.NewInputBox(
		func(t time.Time) {
			displayBox.SetDate(t)
		},
		func(x string) {
			displayBox.SetDescription(x)
	})
	flex := tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputBox, 0, 10, false).
		AddItem(displayBox.GetTextView(), 0, 7, false)
	err := app.SetRoot(flex, true).SetFocus(inputBox).Run()
	if err != nil {
		panic(err)
	}
}
