package app_test

import (
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	logrusTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/accountguesser"
	"github.com/vitorqb/addledger/internal/ammountguesser"
	. "github.com/vitorqb/addledger/internal/app"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/testutils"
	accountguesser_mock "github.com/vitorqb/addledger/mocks/accountguesser"
	ammountguesser_mock "github.com/vitorqb/addledger/mocks/ammountguesser"
	. "github.com/vitorqb/addledger/mocks/dateguesser"
	. "github.com/vitorqb/addledger/mocks/transactionmatcher"
)

var account = journal.Account("ACC")

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
		state.Transaction.Description.Set(description)

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
		statementEntry := finance.StatementEntry{}

		// The expected guess
		expectedGuess := *testutils.Ammount_1(t)

		// Set on state
		state.InputMetadata.SetMatchingTransactions(matchingTransactions)
		state.InputMetadata.SetPostingAmmountText(userInput)
		newPosting := statemod.NewPostingData()
		state.Transaction.Postings.Append(newPosting)
		state.SetStatementEntries([]finance.StatementEntry{statementEntry})
		// Assertions & behavior for engine
		guesser.EXPECT().Guess(gomock.Any()).DoAndReturn(func(inputs ammountguesser.Inputs) (finance.Ammount, bool) {
			assert.Equal(t, userInput, inputs.UserInput)
			assert.Equal(t, statementEntry, inputs.StatementEntry)
			assert.Equal(t, matchingTransactions, inputs.MatchingTransactions)
			assert.Equal(t, state.Transaction.Postings.Get(), inputs.PostingsData)
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

func TestLinkAccountGuesser(t *testing.T) {
	t.Run("Updates guess in state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		guesser := accountguesser_mock.NewMockAccountGuesser(ctrl)

		// Input used for guessing
		statementEntry := testutils.StatementEntry_1(t)
		statamentEntries := []finance.StatementEntry{statementEntry}
		matchingTransactions := []journal.Transaction{*testutils.Transaction_1(t)}
		postings := []journal.Posting{testutils.Posting_1(t)}
		postingData := testutils.PostingData_1(t)
		postingsData := []*statemod.PostingData{&postingData}
		transationHistory := []journal.Transaction{*testutils.Transaction_2(t)}
		userInput := "User input"
		inputs := accountguesser.Inputs{
			MatchingTransactions: matchingTransactions,
			PostingInputs:        postings,
			Description:          userInput,
			TransactionHistory:   transationHistory,
			StatementEntry:       statementEntry,
		}

		// Set on state
		state := statemod.InitialState()
		state.SetStatementEntries(statamentEntries)
		state.InputMetadata.SetMatchingTransactions(matchingTransactions)
		state.Transaction.Postings.Set(postingsData)
		state.Transaction.Description.Set(userInput)
		state.JournalMetadata.SetTransactions(transationHistory)

		// The expected call to guesser
		guesser.EXPECT().Guess(inputs).Return(journal.Account(account), true)

		// Links
		LinkAccountGuesser(state, guesser)

		// Forces state hooks to run
		state.SetPhase(statemod.InputPostingAccount)

		// Ensures state is updated w guess
		actualGuess, _ := state.InputMetadata.GetPostingAccountGuess()
		assert.Equal(t, account, actualGuess)
	})

	t.Run("Clears guess when no guess", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		guesser := accountguesser_mock.NewMockAccountGuesser(ctrl)

		// Set a guess on state
		state := statemod.InitialState()
		state.InputMetadata.SetPostingAccountGuess(account)

		// The expected call to guesser
		guesser.EXPECT().Guess(gomock.Any()).Return(journal.Account(""), false)

		// Links
		LinkAccountGuesser(state, guesser)

		// Forces state hooks to run
		state.SetPhase(statemod.InputPostingAccount)

		// Ensures state is updated w guess
		actualGuess, success := state.InputMetadata.GetPostingAccountGuess()
		assert.Equal(t, journal.Account(""), actualGuess)
		assert.False(t, success)
	})

}

func TestLinkDateGuesser(t *testing.T) {

	aDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	aDateStr := aDate.Format("2006-01-02")
	sEntry := finance.StatementEntry{}

	testcases := []struct {
		name string
		run  func(*testing.T, *statemod.State, *MockIDateGuesser)
	}{
		{
			name: "Updates guess in state from new date input",
			run: func(t *testing.T, state *statemod.State, guesser *MockIDateGuesser) {
				// Expected guesser call
				guesser.EXPECT().Guess(aDateStr, sEntry).Return(aDate, true)

				// Link
				LinkDateGuesser(state, guesser)

				// Set inputs on state
				state.InputMetadata.SetDateText(aDateStr)

				// Ensure guess was set
				actualGuess, _ := state.InputMetadata.GetDateGuess()
				assert.Equal(t, aDate, actualGuess)
			},
		},
		{
			name: "Updates guess in state from new statement entry",
			run: func(t *testing.T, state *statemod.State, guesser *MockIDateGuesser) {
				// Expected guesser call
				guesser.EXPECT().Guess("", sEntry).Return(aDate, true)

				// Link
				LinkDateGuesser(state, guesser)

				// Set inputs on state
				state.SetStatementEntries([]finance.StatementEntry{sEntry})

				// Ensure guess was set
				actualGuess, _ := state.InputMetadata.GetDateGuess()
				assert.Equal(t, aDate, actualGuess)
			},
		},
		{
			name: "Clears guess when no guess",
			run: func(t *testing.T, state *statemod.State, guesser *MockIDateGuesser) {
				// Expected guesser call
				guesser.EXPECT().Guess(aDateStr, sEntry).Return(time.Time{}, false)

				// Starts with a guess
				state.InputMetadata.SetDateGuess(aDate)

				// Link
				LinkDateGuesser(state, guesser)

				// Set inputs on state
				state.InputMetadata.SetDateText(aDateStr)

				// Ensure guess was cleared
				_, found := state.InputMetadata.GetDateGuess()
				assert.False(t, found)
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			guesser := NewMockIDateGuesser(ctrl)
			testcase.run(t, statemod.InitialState(), guesser)
		})
	}
}
