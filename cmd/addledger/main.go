package main

import (
	"os"

	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/controller"
	"github.com/vitorqb/addledger/internal/display"
	"github.com/vitorqb/addledger/internal/state"
)

func main() {

	config, err := config.LoadFromCommandLine()
	if err != nil {
		logrus.WithError(err).Fatal("Error loading config.")
	}

	state := state.InitialState()
	app := tview.NewApplication()

	destFile, err := os.OpenFile(config.DestFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logrus.WithError(err).
			WithField("file", config.DestFile).
			Fatal("Failed to open file")
	}

	controller, err := controller.NewController(state, controller.WithOutput(destFile))
	if err != nil {
		logrus.WithError(err).Fatal("Failed to instantiate controller")
	}
	layout := display.NewLayout(controller, state)
	err = app.
		SetRoot(layout.GetContent(), true).
		SetFocus(layout.Input.GetContent()).
		Run()
	if err != nil {
		panic(err)
	}
}
