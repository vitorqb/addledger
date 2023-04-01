package main

import (
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/displaybox"
	"github.com/vitorqb/addledger/internal/inputbox"
)

func main() {
	app := tview.NewApplication()
	inputBox := inputbox.NewInputBox()
	displayBox := displaybox.NewDisplayBox()
	flex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(inputBox, 0, 10, false).AddItem(displayBox, 0, 7, false)
	err := app.SetRoot(flex, true).Run()
	if err != nil {
		panic(err)
	}
}
