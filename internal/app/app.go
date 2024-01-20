// The app package contains the logic to manage and run the entire application. It
// coordinates most the functionality of the app and the state management.
package app

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/accountguesser"
	"github.com/vitorqb/addledger/internal/ammountguesser"
	"github.com/vitorqb/addledger/internal/dateguesser"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/transactionmatcher"
	"github.com/vitorqb/addledger/internal/userinput"
)

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

		descriptionInput, found := state.Transaction.Description.Get()
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
			PostingsData:         state.Transaction.Postings.Get(),
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
		completePosting := userinput.ExtractPostings(state.Transaction.Postings.Get())
		transactionHist := state.JournalMetadata.Transactions()
		statementEntry, _ := state.CurrentStatementEntry()
		description, _ := state.Transaction.Description.Get()
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

// LinkDateGuesser links a given Date Guesser to state updates
func LinkDateGuesser(state *statemod.State, guesser dateguesser.IDateGuesser) {
	busy := false
	state.AddOnChangeHook(func() {
		if busy {
			return
		}
		busy = true
		defer func() { busy = false }()

		dateText := state.InputMetadata.GetDateText()
		statementEntry, _ := state.CurrentStatementEntry()
		dateGuess, success := guesser.Guess(dateText, statementEntry)
		if success {
			state.InputMetadata.SetDateGuess(dateGuess)
			return
		}
		state.InputMetadata.ClearDateGuess()
	})
}
