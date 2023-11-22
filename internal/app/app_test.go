package app_test

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	logrusTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/ammountguesser"
	. "github.com/vitorqb/addledger/internal/app"
	"github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/statementloader"
	"github.com/vitorqb/addledger/internal/testutils"
	ammountguesser_mock "github.com/vitorqb/addledger/mocks/ammountguesser"
	. "github.com/vitorqb/addledger/mocks/statementloader"
	. "github.com/vitorqb/addledger/mocks/transactionmatcher"
)

func TestLoadStatement(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		loader := NewMockStatementLoader(ctrl)
		statementEntries := []statementloader.StatementEntry{
			{Account: "ACC", Description: "FOO"},
			{Account: "ACC", Description: "BAR"},
		}
		loader.EXPECT().Load(gomock.Any()).Return(statementEntries, nil)
		filePath := testutils.TestDataPath(t, "file")
		state := statemod.InitialState()
		err := LoadStatement(loader, filePath, state)
		assert.Nil(t, err)
		assert.Equal(t, statementEntries, state.GetStatementEntries())
	})

}

func TestConfigureLogger(t *testing.T) {

	t.Run("Logs to a file", func(t *testing.T) {
		logger := logrus.New()
		tempDir := t.TempDir()
		logFile := tempDir + "/log"
		ConfigureLogger(logger, logFile, "info")
		logger.Info("foo")

		file, err := os.Open(logFile)
		assert.Nil(t, err)
		defer file.Close()

		fileContentBytes, err := os.ReadFile(logFile)
		assert.Nil(t, err)
		fileContent := string(fileContentBytes)
		assert.Contains(t, fileContent, "foo")
	})

	t.Run("Fails to open log file", func(t *testing.T) {
		logger, logHook := logrusTest.NewNullLogger()
		logger.ExitFunc = func(int) {}
		logFile := "/foo/bar"
		ConfigureLogger(logger, logFile, "info")
		assert.Equal(t, 1, len(logHook.Entries))
		assert.Equal(t, logrus.FatalLevel, logHook.LastEntry().Level)
		assert.Contains(t, logHook.LastEntry().Message, "Failed to open log file")
	})

	t.Run("Fails to parse log level", func(t *testing.T) {
		logger, logHook := logrusTest.NewNullLogger()
		logger.ExitFunc = func(int) {}
		ConfigureLogger(logger, "", "foo")
		assert.Equal(t, 1, len(logHook.Entries))
		assert.Equal(t, logrus.FatalLevel, logHook.LastEntry().Level)
		assert.Contains(t, logHook.LastEntry().Message, "Failed to parse log level")
	})
}

func TestMaybeLoadCsvStatement(t *testing.T) {
	t.Run("No file", func(t *testing.T) {
		state := statemod.InitialState()
		err := MaybeLoadCsvStatement(config.CSVStatementLoaderConfig{}, state)
		assert.Nil(t, err)
		assert.Equal(t, []statementloader.StatementEntry{}, state.GetStatementEntries())
	})
	t.Run("Fails to load statement", func(t *testing.T) {
		state := statemod.InitialState()
		err := MaybeLoadCsvStatement(config.CSVStatementLoaderConfig{File: "dont-exist"}, state)
		assert.ErrorContains(t, err, "failed to load statement")
	})
}

func TestLinkTransactionMatcher(t *testing.T) {
	t.Run("Saves transactions to state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Prepares state with description and transaction history
		transactionHistory := []journal.Transaction{{Comment: "one"}}
		description := "description"
		state := statemod.InitialState()

		// Prepares the matcher with expected calls
		matchedTransactions := []journal.Transaction{{Comment: "two"}}
		matcher := NewMockITransactionMatcher(ctrl)
		matcher.EXPECT().SetDescriptionInput(description)
		matcher.EXPECT().SetTransactionHistory(transactionHistory)
		matcher.EXPECT().Match().Return(matchedTransactions)

		// Links
		LinkTransactionMatcher(state, matcher)

		// Set the state variables
		state.JournalMetadata.SetTransactions(transactionHistory)
		state.JournalEntryInput.SetDescription(description)

		// Check state was properly set
		resultMatchingTransactions := state.InputMetadata.MatchingTransactions()
		assert.Equal(t, matchedTransactions, resultMatchingTransactions)
	})
}

func TestLinkAmmountGuesser(t *testing.T) {
	t.Run("Saves guess to state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		guesser := ammountguesser_mock.NewMockIAmmountGuesser(ctrl)
		state := statemod.InitialState()

		// All input used for guessing
		matchingTransactions := []journal.Transaction{*testutils.Transaction_1(t)}
		userInput := "EUR 12.20"
		statementEntry := statementloader.StatementEntry{}

		// The expected guess
		expectedGuess := *testutils.Ammount_1(t)

		// Set on state
		state.InputMetadata.SetMatchingTransactions(matchingTransactions)
		state.InputMetadata.SetPostingAmmountText(userInput)
		state.JournalEntryInput.AddPosting()
		state.SetStatementEntries([]statementloader.StatementEntry{statementEntry})
		// Assertions & behavior for engine
		guesser.EXPECT().Guess(gomock.Any()).DoAndReturn(func(inputs ammountguesser.Inputs) (finance.Ammount, bool) {
			assert.Equal(t, userInput, inputs.UserInput)
			assert.Equal(t, statementEntry, inputs.StatementEntry)
			assert.Equal(t, matchingTransactions, inputs.MatchingTransactions)
			assert.Equal(t, state.JournalEntryInput.GetPostings(), inputs.PostingInputs)
			return expectedGuess, true
		})

		// Prepares the engine
		LinkAmmountGuesser(state, guesser)

		// Trigger state change hooks
		state.SetPhase(statemod.InputPostingAmmount)

		// Ensure ammount guesser guessed
		actualGuess, _ := state.InputMetadata.GetPostingAmmountGuess()
		assert.Equal(t, expectedGuess, actualGuess)
	})
}
