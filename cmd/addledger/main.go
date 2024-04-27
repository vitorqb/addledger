package main

import (
	"os"

	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/app"
	"github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/controller"
	"github.com/vitorqb/addledger/internal/display"
	"github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/injector"
	"github.com/vitorqb/addledger/internal/services/statementloader"
)

func main() {

	// Loads config
	config, err := config.LoadFromCommandLine()
	if err != nil {
		logrus.WithError(err).Fatal("Error loading config.")
	}

	// Configures logging
	app.ConfigureLogger(logrus.StandardLogger(), config.LogFile, config.LogLevel)

	// Creates a hledger client
	hledgerClient := injector.HledgerClient(config)

	// Loads state
	state, err := injector.State(hledgerClient)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load state")
	}

	// Loads metadata
	metaLoader, err := injector.MetaLoader(state, hledgerClient)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load metadata loader")
	}
	err = metaLoader.LoadAccounts()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load accounts")
	}
	err = metaLoader.LoadTransactions()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load transactions")
	}

	// Opens the destination file
	destFile, err := os.OpenFile(config.DestFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logrus.WithError(err).
			WithField("file", config.DestFile).
			Fatal("Failed to open file")
	}

	// Starts the EventBus
	eventBus := eventbus.New()

	// Starts a date guesser
	dateGuesser, err := injector.DateGuesser()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load date guesser")
	}
	app.LinkDateGuesser(state, dateGuesser)

	// Starts a Printer
	printer, printerErr := injector.Printer(config.PrinterConfig)
	if printerErr != nil {
		logrus.WithError(err).Fatal("Failed to load printer")
	}

	// Loads a TransactionMatcher. We don't need the reference since it's
	// linked to the state.
	transactionMatcher, err := injector.TransactionMatcher()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load transaction matcher")
	}
	app.LinkTransactionMatcher(state, transactionMatcher)

	// Prepares a statement loader
	statementReader := injector.StatementReader()
	statementLoaderSvc := statementloader.New(state, statementReader)

	// Starts a user messenger
	userMessenger := injector.UserMessenger(state)

	// Starts a new controller
	controller, err := controller.NewController(state,
		controller.WithOutput(destFile),
		controller.WithEventBus(eventBus),
		controller.WithDateGuesser(dateGuesser),
		controller.WithMetaLoader(metaLoader),
		controller.WithPrinter(printer),
		controller.WithCSVStatementLoader(statementLoaderSvc),
		controller.WithUserMessenger(userMessenger),
	)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to instantiate controller")
	}

	// Starts the AmmountGuesserEngine. Note it's linked to state refresh
	// so we don't need it's instance.
	ammountGuesser := injector.AmmountGuesser()
	app.LinkAmmountGuesser(state, ammountGuesser)

	// Start an account guesser
	accountGuesser, err := injector.AccountGuesser(state)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load account guesser")
	}
	app.LinkAccountGuesser(state, accountGuesser)

	// Maybe load a CSV statement
	if config.CSVStatementFile != "" {
		err = statementLoaderSvc.LoadFromFiles(config.CSVStatementFile, config.CSVStatementPreset)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to load csv statement")
		}
	}

	// Starts a new tview App
	app := tview.NewApplication()

	// Starts a new layout
	layout, err := display.NewLayout(controller, state, eventBus, app)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to instatiate layout")
	}

	// Run!
	err = app.SetRoot(layout, true).SetFocus(layout).Run()
	if err != nil {
		logrus.Fatal(err)
	}
}
