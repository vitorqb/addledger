package main

import (
	"time"

	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/displaybox"
	"github.com/vitorqb/addledger/internal/inputbox"
	"github.com/vitorqb/addledger/internal/inputcontextbox"
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
	inputContextBox := inputcontextbox.NewInputContextBox()
	flex := tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(displayBox.GetTextView(), 0, 5, false).
		AddItem(inputBox.GetInputField(), 0, 1, false).
		AddItem(inputContextBox.GetTextView(), 0, 10, false)
	err := app.SetRoot(flex, true).SetFocus(inputBox.GetInputField()).Run()
	if err != nil {
		panic(err)
	}
}
