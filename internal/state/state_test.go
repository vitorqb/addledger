package state_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
	. "github.com/vitorqb/addledger/internal/state"
)

var anAmmount = finance.Ammount{
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
			name: "Notify on change of Display",
			run: func(t *testing.T, c *testcontext) {
				c.state.Display.SetShortcutModal(false)
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
				assert.Equal(t, c.state.CurrentPhase(), InputTags)
				assert.Equal(t, 2, c.hookCallCounter)

				c.state.NextPhase()
				assert.Equal(t, c.state.CurrentPhase(), InputPostingAccount)
				assert.Equal(t, 3, c.hookCallCounter)
			},
		},
		{
			name: "PrevPhase",
			run: func(t *testing.T, c *testcontext) {
				c.state.SetPhase(Confirmation)
				assert.Equal(t, c.state.CurrentPhase(), Confirmation)

				c.state.PrevPhase()
				assert.Equal(t, c.state.CurrentPhase(), InputPostingAmmount)
				assert.Equal(t, 2, c.hookCallCounter)

				c.state.PrevPhase()
				assert.Equal(t, c.state.CurrentPhase(), InputPostingAccount)
				assert.Equal(t, 3, c.hookCallCounter)

				c.state.PrevPhase()
				assert.Equal(t, c.state.CurrentPhase(), InputTags)
				assert.Equal(t, 4, c.hookCallCounter)
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
		{
			name: "Manipulates selected tag",
			run: func(t *testing.T, c *testcontext) {
				tag := journal.Tag{Name: "FOO", Value: "BAR"}
				assert.Empty(t, c.state.InputMetadata.SelectedTag())
				c.state.InputMetadata.SetSelectedTag(tag)
				assert.Equal(t, tag, c.state.InputMetadata.SelectedTag())
				assert.Equal(t, 1, c.hookCallCounter)
			},
		},
		{
			name: "InputMetadata resets properly",
			run: func(t *testing.T, c *testcontext) {
				c.state.InputMetadata.SetPostingAccountText("FOO")
				c.state.InputMetadata.SetSelectedPostingAccount("BAR")
				c.state.InputMetadata.SetDescriptionText("FOO")
				c.state.InputMetadata.SetSelectedDescription("BAR")
				c.state.InputMetadata.SetPostingAmmountGuess(anAmmount)
				c.state.InputMetadata.SetPostingAmmountInput(anAmmount)
				c.state.InputMetadata.SetPostingAmmountText("FOO")
				c.state.InputMetadata.SetDateGuess(time.Time{})
				assert.Equal(t, 8, c.hookCallCounter)

				c.state.InputMetadata.Reset()
				assert.Equal(t, 9, c.hookCallCounter)

				postingAccountText := c.state.InputMetadata.PostingAccountText()
				assert.Equal(t, "", postingAccountText)
				selectedPostingAccount := c.state.InputMetadata.SelectedPostingAccount()
				assert.Equal(t, "", selectedPostingAccount)
				descriptionText := c.state.InputMetadata.DescriptionText()
				assert.Equal(t, "", descriptionText)
				selectedDescription := c.state.InputMetadata.SelectedDescription()
				assert.Equal(t, "", selectedDescription)
				_, found := c.state.InputMetadata.GetPostingAmmountGuess()
				assert.False(t, found)
				_, found = c.state.InputMetadata.GetPostingAmmountInput()
				assert.False(t, found)
				postingAmmountText := c.state.InputMetadata.GetPostingAmmountText()
				assert.Equal(t, "", postingAmmountText)
				_, found = c.state.InputMetadata.GetDateGuess()
				assert.False(t, found)
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

func TestJournalMetadata(t *testing.T) {
	type testcontext struct {
		hookCallCounter int
		journalMetadata *JournalMetadata
	}

	type testcase struct {
		name string
		run  func(t *testing.T, c *testcontext)
	}

	testcases := []testcase{
		{
			name: "Manipulate accounts",
			run: func(t *testing.T, c *testcontext) {
				assert.Empty(t, c.journalMetadata.Accounts())
				accs := []journal.Account{"FOO", "BAR"}
				c.journalMetadata.SetAccounts(accs)
				assert.Equal(t, accs, c.journalMetadata.Accounts())
				assert.Equal(t, 1, c.hookCallCounter)
			},
		},
		{
			name: "Manipulate tags",
			run: func(t *testing.T, c *testcontext) {
				assert.Empty(t, c.journalMetadata.Tags())
				tags := []journal.Tag{{Name: "FOO", Value: "BAR"}}
				c.journalMetadata.SetTransactions([]journal.Transaction{{Tags: tags}})
				assert.Equal(t, tags, c.journalMetadata.Tags())
				assert.Equal(t, 1, c.hookCallCounter)
			},
		},
		{
			name: "Remove duplicat tags",
			run: func(t *testing.T, c *testcontext) {
				assert.Empty(t, c.journalMetadata.Tags())
				tags := []journal.Tag{{Name: "FOO", Value: "BAR"}}
				transaction1 := journal.Transaction{Tags: tags}
				transaction2 := journal.Transaction{Tags: tags}
				transactions := []journal.Transaction{transaction1, transaction2}
				c.journalMetadata.SetTransactions(transactions)
				assert.Equal(t, tags, c.journalMetadata.Tags())
				assert.Equal(t, 1, c.hookCallCounter)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := new(testcontext)
			c.hookCallCounter = 0
			c.journalMetadata = NewJournalMetadata()
			c.journalMetadata.AddOnChangeHook(func() { c.hookCallCounter++ })
			tc.run(t, c)
		})
	}

}

func TestDisplay(t *testing.T) {
	{
		type testcontext struct {
			hookCallCounter int
			display         *Display
		}

		type testcase struct {
			name string
			run  func(t *testing.T, c *testcontext)
		}

		testcases := []testcase{
			{
				name: "Manipulate ShortcutModal",
				run: func(t *testing.T, c *testcontext) {
					assert.False(t, c.display.ShortcutModal())
					c.display.SetShortcutModal(true)
					assert.True(t, c.display.ShortcutModal())
					assert.Equal(t, 1, c.hookCallCounter)
				},
			},
		}

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				c := new(testcontext)
				c.hookCallCounter = 0
				c.display = NewDisplay()
				c.display.AddOnChangeHook(func() { c.hookCallCounter++ })
				tc.run(t, c)
			})
		}

	}
}
