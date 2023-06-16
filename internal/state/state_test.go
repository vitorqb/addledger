package state_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/journal"
	. "github.com/vitorqb/addledger/internal/state"
	hledger_mock "github.com/vitorqb/addledger/mocks/hledger"
)

var anAmmount = journal.Ammount{
	Commodity: "EUR",
	Quantity:  decimal.New(2400, -2),
}

func TestState(t *testing.T) {

	type testcontext struct {
		hookCallCounter int
		state           *State
	}

	type testcase struct {
		name string
		run  func(t *testing.T, c *testcontext)
	}

	testcases := []testcase{
		{
			name: "Notify on change of JournalEntryInput",
			run: func(t *testing.T, c *testcontext) {
				c.state.JournalEntryInput.SetDescription("FOO")
				assert.Equal(t, 1, c.hookCallCounter)

			},
		},
		{
			name: "NextPhase",
			run: func(t *testing.T, c *testcontext) {
				assert.Equal(t, c.state.CurrentPhase(), InputDate)
				c.state.NextPhase()
				assert.Equal(t, c.state.CurrentPhase(), InputDescription)
				assert.Equal(t, 1, c.hookCallCounter)
				c.state.NextPhase()
				assert.Equal(t, c.state.CurrentPhase(), InputPostingAccount)
				assert.Equal(t, 2, c.hookCallCounter)

			},
		},
		{
			name: "SetPhase",
			run: func(t *testing.T, c *testcontext) {
				c.state.SetPhase(InputPostingAccount)
				assert.Equal(t, InputPostingAccount, c.state.CurrentPhase())
				assert.Equal(t, 1, c.hookCallCounter)

				c.state.SetPhase(InputDescription)
				assert.Equal(t, InputDescription, c.state.CurrentPhase())
				assert.Equal(t, 2, c.hookCallCounter)
			},
		},
		{
			name: "SetAccounts",
			run: func(t *testing.T, c *testcontext) {
				c.state.SetAccounts([]string{"FOO"})
				assert.Equal(t, 1, c.hookCallCounter)
				accounts := c.state.GetAccounts()
				assert.Equal(t, []string{"FOO"}, accounts)
			},
		},
		{
			name: "LoadMetadata",
			run: func(t *testing.T, c *testcontext) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				transactions := []journal.Transaction{{Description: "FOO"}, {Description: "Bar"}}
				hledgerClient := hledger_mock.NewMockIClient(ctrl)
				hledgerClient.EXPECT().Accounts().Return([]string{"FOO"}, nil)
				hledgerClient.EXPECT().Transactions().Return(transactions, nil)

				err := c.state.LoadMetadata(hledgerClient)
				assert.Nil(t, err)
				assert.Equal(t, 2, c.hookCallCounter)
				assert.Equal(t, []string{"FOO"}, c.state.GetAccounts())
				assert.Equal(t, transactions, c.state.JournalMetadata.Transactions())
			},
		},
		{
			name: "Manipulates SelectedPostingAccount",
			run: func(t *testing.T, c *testcontext) {
				c.state.InputMetadata.SetSelectedPostingAccount("FOO")
				assert.Equal(t, 1, c.hookCallCounter)
				acc := c.state.InputMetadata.SelectedPostingAccount()
				assert.Equal(t, "FOO", acc)
			},
		},
		{
			name: "Manipulates PostingAmmountGuess",
			run: func(t *testing.T, c *testcontext) {
				_, found := c.state.InputMetadata.GetPostingAmmountGuess()
				assert.False(t, found)
				c.state.InputMetadata.SetPostingAmmountGuess(anAmmount)
				assert.Equal(t, 1, c.hookCallCounter)
				ammount, found := c.state.InputMetadata.GetPostingAmmountGuess()
				assert.True(t, found)
				assert.Equal(t, anAmmount, ammount)
				c.state.InputMetadata.ClearPostingAmmountGuess()
				_, found = c.state.InputMetadata.GetPostingAmmountGuess()
				assert.False(t, found)
				assert.Equal(t, 2, c.hookCallCounter)
			},
		},
		{
			name: "Manipulates PostingAmmountInput",
			run: func(t *testing.T, c *testcontext) {
				_, found := c.state.InputMetadata.GetPostingAmmountInput()
				assert.False(t, found)
				c.state.InputMetadata.SetPostingAmmountInput(anAmmount)
				assert.Equal(t, 1, c.hookCallCounter)
				ammount, found := c.state.InputMetadata.GetPostingAmmountInput()
				assert.True(t, found)
				assert.Equal(t, anAmmount, ammount)
				c.state.InputMetadata.ClearPostingAmmountInput()
				_, found = c.state.InputMetadata.GetPostingAmmountInput()
				assert.False(t, found)
				assert.Equal(t, 2, c.hookCallCounter)
			},
		},
		{
			name: "Manipulates PostingAmmountText",
			run: func(t *testing.T, c *testcontext) {
				text := c.state.InputMetadata.GetPostingAmmountText()
				assert.Equal(t, "", text)
				c.state.InputMetadata.SetPostingAmmountText("EUR 12.20")
				assert.Equal(t, 1, c.hookCallCounter)
				text = c.state.InputMetadata.GetPostingAmmountText()
				assert.Equal(t, "EUR 12.20", text)
				c.state.InputMetadata.ClearPostingAmmountText()
				text = c.state.InputMetadata.GetPostingAmmountText()
				assert.Equal(t, "", text)
				assert.Equal(t, 2, c.hookCallCounter)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := new(testcontext)
			c.hookCallCounter = 0
			c.state = InitialState()
			c.state.AddOnChangeHook(func() { c.hookCallCounter++ })
			tc.run(t, c)
		})
	}
}
