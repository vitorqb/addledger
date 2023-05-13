package main

import (
	"os"

	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/controller"
	"github.com/vitorqb/addledger/internal/display"
	"github.com/vitorqb/addledger/internal/injector"
	"github.com/vitorqb/addledger/internal/state"
)

func main() {

	// Loads config
	config, err := config.LoadFromCommandLine()
	if err != nil {
		logrus.WithError(err).Fatal("Error loading config.")
	}

	// Configures logging
	if config.LogFile != "" {
		logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to open log file.")
		}
		logrus.SetOutput(logFile)
	}
	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to parse log level.")
	}
	logrus.SetLevel(logLevel)

	// Creates a hledger client
	hledgerClient := injector.HledgerClient(config)

	// Loads state w/ metadata
	state := state.InitialState()
	err = state.LoadMetadata(hledgerClient)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load metadata")
	}

	// Opens the destination file
	destFile, err := os.OpenFile(config.DestFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logrus.WithError(err).
			WithField("file", config.DestFile).
			Fatal("Failed to open file")
	}

	// Starts a new controller
	controller, err := controller.NewController(state, controller.WithOutput(destFile))
	if err != nil {
		logrus.WithError(err).Fatal("Failed to instantiate controller")
	}

	// Starts a new layout
	layout := display.NewLayout(controller, state)

	// Starts a new tview App
	app := tview.NewApplication()

	// Run!
	err = app.
		SetRoot(layout.GetContent(), true).
		SetFocus(layout.Input.GetContent()).
		Run()
	if err != nil {
		panic(err)
	}
}
