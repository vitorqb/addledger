package ammountguesser_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/ammountguesser"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/state"
	tu "github.com/vitorqb/addledger/internal/testutils"
)

var anAmmount = finance.Ammount{Commodity: "EUR", Quantity: decimal.New(1221, -2)}
var anAmmountBRL = finance.Ammount{Commodity: "BRL", Quantity: decimal.New(1222, -2)}
var anotherAmmount = finance.Ammount{Commodity: "USD", Quantity: decimal.New(9922, -2)}

func TestAmmountGuesser(t *testing.T) {
	type testcase struct {
		name      string
		inputs    Inputs
		setupFunc func(tc *testcase)
		guess     finance.Ammount
		success   bool
	}
	var testcases = []testcase{
		{
			name:    "Guesses from user input when valid stirng",
			inputs:  Inputs{UserInput: "EUR 12.21"},
			guess:   anAmmount,
			success: true,
		},
		{
			name:    "Guesses from user input without currency using default",
			inputs:  Inputs{UserInput: "12.21"},
			guess:   anAmmount,
			success: true,
		},
		{
			name:    "Guesses from user input w another currency",
			inputs:  Inputs{UserInput: "BRL 12.22"},
			guess:   anAmmountBRL,
			success: true,
		},
		{
			name: "Guesses from matching transaction",
			setupFunc: func(tc *testcase) {
				t := tu.Transaction_2(t)
				tc.inputs.MatchingTransactions = []journal.Transaction{*t}
				tc.guess = t.Posting[0].Ammount
			},
			success: true,
		},
		{
			name: "Don't guess from matching transaction if user input text",
			setupFunc: func(tc *testcase) {
				t := tu.Transaction_1(t)
				tc.inputs.MatchingTransactions = []journal.Transaction{*t}
				tc.inputs.UserInput = "EUR 12.21"
			},
			guess:   anAmmount,
			success: true,
		},
		{
			name: "Guess from loaded statement entry",
			setupFunc: func(tc *testcase) {
				// Set some matching transactions that should be ignored
				t := tu.Transaction_2(t)
				tc.inputs.MatchingTransactions = []journal.Transaction{*t}

				// Set a statement entry
				tc.inputs.StatementEntry = finance.StatementEntry{Ammount: anotherAmmount}

			},
			guess:   anotherAmmount.InvertSign(),
			success: true,
		},
		{
			name: "Guess from pending balance",
			setupFunc: func(tc *testcase) {
				// Set some matching transactions that should be ignored
				transact := tu.Transaction_2(t)
				tc.inputs.MatchingTransactions = []journal.Transaction{*transact}

				// Set a statement entry that should be ignored
				tc.inputs.StatementEntry = finance.StatementEntry{Ammount: anotherAmmount}

				// Set some pending balance
				postingData := state.NewPostingData()
				tu.FillPostingData_1(t, postingData)
				postingData.Ammount.Set(anAmmount)
				tc.inputs.PostingsData = []*state.PostingData{postingData}
			},
			guess:   anAmmount.InvertSign(),
			success: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupFunc != nil {
				tc.setupFunc(&tc)
			}
			guesser := AmmountGuesser{}

			guess, success := guesser.Guess(tc.inputs)
			assert.Equal(t, tc.guess, guess)
			assert.Equal(t, tc.success, success)
		})
	}
}
