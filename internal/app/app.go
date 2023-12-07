// The app package contains the logic to manage and run the entire application. It
// coordinates most the functionality of the app and the state management.
package app

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/accountguesser"
	"github.com/vitorqb/addledger/internal/ammountguesser"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/statementreader"
	"github.com/vitorqb/addledger/internal/transactionmatcher"
)

// LoadStatement loads a statement from a file and saves it to the state.
func LoadStatement(
	loader statementreader.StatementReader,
	file string,
	state *statemod.State,
) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	entries, err := loader.Read(f)
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

// LinkTransactionMatcher links a transaction matcher to the state. Every time the
// state changes, the matcher will calculate the transactions matching the user
// inputs and save them to the state.
func LinkTransactionMatcher(state *statemod.State, matcher transactionmatcher.ITransactionMatcher) {
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

// LinkAmmountGuesser links a given Ammount Guesser to state updates
func LinkAmmountGuesser(state *statemod.State, guesser ammountguesser.IAmmountGuesser) {
	busy := false
	// subscribes to changes
	state.AddOnChangeHook(func() {
		if busy {
			return
		}
		busy = true
		defer func() { busy = false }()

		currentStatementEntry, _ := state.CurrentStatementEntry()
		inputs := ammountguesser.Inputs{
			UserInput:            state.InputMetadata.GetPostingAmmountText(),
			PostingInputs:        state.JournalEntryInput.GetPostings(),
			StatementEntry:       currentStatementEntry,
			MatchingTransactions: state.InputMetadata.MatchingTransactions(),
		}
		guess, success := guesser.Guess(inputs)
		if !success {
			state.InputMetadata.ClearPostingAmmountGuess()
			return
		}
		state.InputMetadata.SetPostingAmmountGuess(guess)
	})
}

// LinkAccountGuesser links a given Account Guesser to state updates
func LinkAccountGuesser(state *statemod.State, guesser accountguesser.AccountGuesser) {
	busy := false
	state.AddOnChangeHook(func() {
		if busy {
			return
		}
		busy = true
		defer func() { busy = false }()

		matchedTransactions := state.InputMetadata.MatchingTransactions()
		completePosting := state.JournalEntryInput.GetCompletePostings()
		transactionHist := state.JournalMetadata.Transactions()
		statementEntry, _ := state.CurrentStatementEntry()
		description, _ := state.JournalEntryInput.GetDescription()
		inputs := accountguesser.Inputs{
			StatementEntry:       statementEntry,
			MatchingTransactions: matchedTransactions,
			Description:          description,
			PostingInputs:        completePosting,
			TransactionHistory:   transactionHist,
		}
		acc, success := guesser.Guess(inputs)
		if !success {
			state.InputMetadata.ClearPostingAccountGuess()
			return
		}
		state.InputMetadata.SetPostingAccountGuess(acc)
	})
}
