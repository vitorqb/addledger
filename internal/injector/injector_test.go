package injector_test

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/ammountguesser"
	"github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/injector"
	. "github.com/vitorqb/addledger/internal/injector"
	"github.com/vitorqb/addledger/internal/journal"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/testutils"
	hledger_mock "github.com/vitorqb/addledger/mocks/hledger"
)

func TestAmmountGuesserEngine(t *testing.T) {
	state := statemod.InitialState()
	_ = AmmountGuesserEngine(state)

	// At the beggining, default guess
	guess, found := state.InputMetadata.GetPostingAmmountGuess()
	assert.True(t, found)
	assert.Equal(t, ammountguesser.DefaultGuess, guess)

	// On new input for ammount guesser text, updates guess
	state.InputMetadata.SetPostingAmmountText("99.99")
	guess, found = state.InputMetadata.GetPostingAmmountGuess()
	assert.True(t, found)
	expectedGuess := finance.Ammount{
		Commodity: ammountguesser.DefaultCommodity,
		Quantity:  decimal.New(9999, -2),
	}
	assert.Equal(t, expectedGuess, guess)

	// On invalid input, defaults to default guess
	state.InputMetadata.SetPostingAmmountText("aaaa")
	guess, found = state.InputMetadata.GetPostingAmmountGuess()
	assert.True(t, found)
	assert.Equal(t, ammountguesser.DefaultGuess, guess)
}

func TestStateAndMetaLoader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	transactions := []journal.Transaction{{Description: "FOO"}, {Description: "Bar"}}
	hledgerClient := hledger_mock.NewMockIClient(ctrl)
	hledgerClient.EXPECT().Accounts().Return([]journal.Account{"FOO"}, nil)
	hledgerClient.EXPECT().Transactions().Return(transactions, nil)

	state, err := State(hledgerClient)
	assert.Nil(t, err)

	metaLoader, err := MetaLoader(state, hledgerClient)
	assert.Nil(t, err)
	err = metaLoader.LoadAccounts()
	assert.Nil(t, err)
	err = metaLoader.LoadTransactions()
	assert.Nil(t, err)
	assert.Equal(t, []journal.Account{"FOO"}, state.JournalMetadata.Accounts())
	assert.Equal(t, transactions, state.JournalMetadata.Transactions())
}

func TestDescriptionMatchAccountGuesser(t *testing.T) {
	state := statemod.InitialState()
	accountGuesser, err := DescriptionMatchAccountGuesser(state)
	if err != nil {
		t.Fatal(err)
	}

	// At the beggining, no guess
	_, success := accountGuesser.Guess()
	assert.False(t, success)

	// Add a TransactionMatcher to state
	_, err = TransactionMatcher(state)
	if err != nil {
		t.Fatal(err)
	}

	// Set the transaction history on state
	state.JournalMetadata.SetTransactions([]journal.Transaction{
		{
			Description: "Supermarket",
			Posting: []journal.Posting{
				{
					Account: "bank:currentaccount",
				},
				{
					Account: "expenses:supermarket",
				},
			},
		},
	})

	// Set a user inputted posting on state
	posting := state.JournalEntryInput.AddPosting()
	posting.SetAccount("foo")
	posting.SetAmmount(finance.Ammount{})

	// Set a user inputted description on state
	state.JournalEntryInput.SetDescription("Superm")

	// Guess should be right
	guess, success := accountGuesser.Guess()
	assert.True(t, success)
	assert.Equal(t, journal.Account("expenses:supermarket"), guess)
}

func TestLastTransactionAccountGuesser(t *testing.T) {
	state := statemod.InitialState()
	accountGuesser, err := LastTransactionAccountGuesser(state)
	if err != nil {
		t.Fatal(err)
	}

	// At the beggining, no guess
	_, success := accountGuesser.Guess()
	assert.False(t, success)

	// Set the transaction history on state
	state.JournalMetadata.SetTransactions([]journal.Transaction{
		{
			Description: "Supermarket",
			Posting: []journal.Posting{
				{
					Account: "bank:currentaccount",
				},
				{
					Account: "expenses:supermarket",
				},
			},
		},
	})

	// Guess should be right
	guess, success := accountGuesser.Guess()
	assert.True(t, success)
	assert.Equal(t, journal.Account("bank:currentaccount"), guess)
}

func TestPrinter(t *testing.T) {
	config := config.PrinterConfig{NumLineBreaksBefore: 2, NumLineBreaksAfter: 3}
	printer, err := injector.Printer(config)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	printer.Print(&buf, *testutils.Transaction_1(t))
	expectedPrint := "\n\n1993-11-23 Description1\n    ACC1    EUR 12.2\n    ACC2    EUR -12.2\n\n\n"
	assert.Equal(t, expectedPrint, buf.String())
}

func TestTransactionMatcher(t *testing.T) {
	state := statemod.InitialState()
	_, err := injector.TransactionMatcher(state)
	assert.Nil(t, err)
	transactions := []journal.Transaction{{Description: "test"}, {Description: "INVALID"}}
	expectedMatchedTransactions := []journal.Transaction{{Description: "test"}}

	// Updates the state
	state.JournalMetadata.SetTransactions(transactions)
	state.JournalEntryInput.SetDescription("test")

	// Ensure that the matched transactions are on the state.
	assert.Equal(t, expectedMatchedTransactions, state.InputMetadata.MatchingTransactions())
}
