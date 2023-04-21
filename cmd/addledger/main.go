package main

import (
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/controller"
	"github.com/vitorqb/addledger/internal/display"
	"github.com/vitorqb/addledger/internal/state"
)

func main() {

	config, err := config.LoadFromCommandLine()
	print(config)
	if err != nil {
		logrus.WithError(err).Fatal("Error loading config.")
	}

	state := state.InitialState()
	app := tview.NewApplication()
	controller := controller.NewController(state)
	layout := display.NewLayout(controller, state)
	err = app.
		SetRoot(layout.GetContent(), true).
		SetFocus(layout.Input.GetContent()).
		Run()
	if err != nil {
		panic(err)
	}
}
