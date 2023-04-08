package main

import (
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/displaybox"
	"github.com/vitorqb/addledger/internal/inputbox"
	"github.com/vitorqb/addledger/internal/inputcontextbox"
	"github.com/vitorqb/addledger/internal/state"
)

func main() {
	state := state.InitialState()
	app := tview.NewApplication()
	displayBox := displaybox.NewDisplayBox(state)
	inputBox := inputbox.NewInputBox(state)
	inputContextBox := inputcontextbox.NewInputContextBox()
	flex := tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(displayBox.GetTextView(), 0, 5, false).
		AddItem(inputBox.GetContent(), 0, 1, false).
		AddItem(inputContextBox.GetTextView(), 0, 10, false)
	err := app.SetRoot(flex, true).SetFocus(inputBox.GetContent()).Run()
	if err != nil {
		panic(err)
	}
}
