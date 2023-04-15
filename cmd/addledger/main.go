package main

import (
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/display"
	"github.com/vitorqb/addledger/internal/state"
)

func main() {
	state := state.InitialState()
	app := tview.NewApplication()
	layout := display.NewLayout(state)
	err := app.
		SetRoot(layout.GetContent(), true).
		SetFocus(layout.Input.GetContent()).
		Run()
	if err != nil {
		panic(err)
	}
}
