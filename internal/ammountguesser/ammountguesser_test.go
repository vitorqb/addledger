package ammountguesser_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/ammountguesser"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/statementloader"
	tu "github.com/vitorqb/addledger/internal/testutils"
)

var anAmmount = finance.Ammount{Commodity: "EUR", Quantity: decimal.New(1221, -2)}
var anAmmountBRL = finance.Ammount{Commodity: "BRL", Quantity: decimal.New(1222, -2)}
var anotherAmmount = finance.Ammount{Commodity: "USD", Quantity: decimal.New(9922, -2)}

func TestEngine(t *testing.T) {
	type testcontext struct {
		engine *Engine
	}
	type testcase struct {
		name string
		run  func(c *testcontext, t *testing.T)
	}
	var testcases = []testcase{
		{
			name: "Guesses from user input when valid stirng",
			run: func(c *testcontext, t *testing.T) {
				c.engine.SetUserInputText("EUR 12.21")
				guess, success := c.engine.Guess()
				assert.True(t, success)
				assert.Equal(t, anAmmount, guess)
			},
		},
		{
			name: "Guesses from user input without currency using default",
			run: func(c *testcontext, t *testing.T) {
				c.engine.SetUserInputText("12.21")
				guess, success := c.engine.Guess()
				assert.True(t, success)
				assert.Equal(t, anAmmount, guess)
			},
		},
		{
			name: "Guesses from user input w another currency",
			run: func(c *testcontext, t *testing.T) {
				c.engine.SetUserInputText("BRL 12.22")
				guess, success := c.engine.Guess()
				assert.True(t, success)
				assert.Equal(t, anAmmountBRL, guess)
			},
		},
		{
			name: "Guesses from matching transaction",
			run: func(c *testcontext, t *testing.T) {
				c.engine.SetUserInputText("")
				transaction := tu.Transaction_2(t)
				matchingTransaction := []journal.Transaction{*transaction}
				c.engine.SetMatchingTransactions(matchingTransaction)
				guess, success := c.engine.Guess()
				assert.True(t, success)
				expected := transaction.Posting[0].Ammount
				assert.Equal(t, expected, guess)
			},
		},
		{
			name: "Don't guess from matching transaction if user input text",
			run: func(c *testcontext, t *testing.T) {
				c.engine.SetUserInputText("EUR 12.21")
				transaction := tu.Transaction_1(t)
				matchingTransaction := []journal.Transaction{*transaction}
				c.engine.SetMatchingTransactions(matchingTransaction)
				guess, success := c.engine.Guess()
				assert.True(t, success)
				assert.Equal(t, anAmmount, guess)
			},
		},
		{
			name: "Guess from loaded statement entry",
			run: func(c *testcontext, t *testing.T) {
				// Set some matching transactions that should be ignored
				c.engine.SetUserInputText("")
				transaction := tu.Transaction_2(t)
				matchingTransaction := []journal.Transaction{*transaction}
				c.engine.SetMatchingTransactions(matchingTransaction)

				// Set a statement entry
				sEntry := statementloader.StatementEntry{Ammount: anotherAmmount}
				c.engine.SetStatementEntry(sEntry)

				// Guess
				guess, success := c.engine.Guess()
				assert.True(t, success)
				assert.Equal(t, anotherAmmount.InvertSign(), guess)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := new(testcontext)
			c.engine = NewEngine()
			tc.run(c, t)
		})
	}
}
