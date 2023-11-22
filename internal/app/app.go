// The app package contains the logic to manage and run the entire application. It
// coordinates most the functionality of the app and the state management.
package app

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/injector"
	"github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/statementloader"
	"github.com/vitorqb/addledger/internal/transactionmatcher"
)

// LoadStatement loads a statement from a file and saves it to the state.
func LoadStatement(
	loader statementloader.StatementLoader,
	file string,
	state *state.State,
) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	entries, err := loader.Load(f)
	if err != nil {
		return fmt.Errorf("failed to load statement: %w", err)
	}

	state.SetStatementEntries(entries)
	return nil
}

// ConfigureLogger configures the logger.
func ConfigureLogger(logger *logrus.Logger, LogFile string, LogLevel string) {
	if LogFile != "" {
		logFile, err := os.OpenFile(LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			logger.WithError(err).Fatal("Failed to open log file.")
		}
		logger.SetOutput(logFile)
	}
	logLevel, err := logrus.ParseLevel(LogLevel)
	if err != nil {
		logger.WithError(err).Fatal("Failed to parse log level.")
	}
	logger.SetLevel(logLevel)
}

// MaybeLoadStatement loads a CSV statement if the config is set.
func MaybeLoadCsvStatement(config config.CSVStatementLoaderConfig, state *state.State) error {
	if config.File == "" {
		return nil
	}
	csvStatementLoader, err := injector.CSVStatementLoader(config)
	if err != nil {
		return fmt.Errorf("failed to load csv statement loader: %w", err)
	}
	err = LoadStatement(csvStatementLoader, config.File, state)
	if err != nil {
		return fmt.Errorf("failed to load statement: %w", err)
	}
	return nil
}

// LinkTransactionMatcher links a transaction matcher to the state. Every time the
// state changes, the matcher will calculate the transactions matching the user
// inputs and save them to the state.
func LinkTransactionMatcher(state *state.State, matcher transactionmatcher.ITransactionMatcher) {
	busy := false
	state.AddOnChangeHook(func() {
		if busy {
			return
		}
		busy = true
		defer func() {
			busy = false
		}()

		descriptionInput, found := state.JournalEntryInput.GetDescription()
		if !found {
			return
		}
		matcher.SetDescriptionInput(descriptionInput)

		transactionHistory := state.JournalMetadata.Transactions()
		matcher.SetTransactionHistory(transactionHistory)

		matchingTransactions := matcher.Match()
		state.InputMetadata.SetMatchingTransactions(matchingTransactions)
	})
}
