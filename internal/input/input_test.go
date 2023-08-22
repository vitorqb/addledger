package input_test

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/internal/journal"
	tu "github.com/vitorqb/addledger/internal/testutils"
)

func TestJournalEntryInput(t *testing.T) {

	type context struct {
		onChangeCalled    bool
		onChangeCallCount int
		input             *JournalEntryInput
	}

	type test struct {
		name string
		run  func(*testing.T, *context)
	}

	tests := []test{
		{
			name: "Date",
			run: func(t *testing.T, c *context) {
				_, found := c.input.GetDate()
				assert.False(t, found)
				c.input.SetDate(tu.Date1(t))
				date, found := c.input.GetDate()
				assert.True(t, found)
				assert.Equal(t, date, tu.Date1(t))
				assert.Equal(t, 1, c.onChangeCallCount)
				c.input.ClearDate()
				_, found = c.input.GetDate()
				assert.False(t, found)
				assert.Equal(t, 2, c.onChangeCallCount)
			},
		},
		{
			name: "Description",
			run: func(t *testing.T, c *context) {
				_, found := c.input.GetDescription()
				assert.False(t, found)
				c.input.SetDescription("FOO")
				description, found := c.input.GetDescription()
				assert.True(t, found)
				assert.Equal(t, description, "FOO")
				assert.Equal(t, 1, c.onChangeCallCount)
				c.input.ClearDescription()
				_, found = c.input.GetDescription()
				assert.False(t, found)
				assert.Equal(t, 2, c.onChangeCallCount)
			},
		},
		{
			name: "Add posting account",
			run: func(t *testing.T, c *context) {
				_, found := c.input.GetPosting(0)
				assert.False(t, found)

				addedPosting := c.input.AddPosting()
				foundPosting, found := c.input.GetPosting(0)
				assert.True(t, found)
				assert.Equal(t, foundPosting, addedPosting)

				_, found = addedPosting.GetAccount()
				assert.False(t, found)

				addedPosting.SetAccount("FOO")
				account, found := addedPosting.GetAccount()
				assert.True(t, found)
				assert.Equal(t, account, "FOO")
				assert.Equal(t, 2, c.onChangeCallCount)
			},
		},
		{
			"Count postings",
			func(t *testing.T, c *context) {
				assert.Equal(t, 0, c.input.CountPostings())
				c.input.AddPosting()
				assert.Equal(t, 1, c.input.CountPostings())
				c.input.AddPosting()
				assert.Equal(t, 2, c.input.CountPostings())
			},
		},
		{
			"Delete last posting",
			func(t *testing.T, c *context) {
				assert.Equal(t, 0, c.input.CountPostings())
				assert.Equal(t, 0, c.onChangeCallCount)
				c.input.AddPosting()
				assert.Equal(t, 1, c.input.CountPostings())
				assert.Equal(t, 1, c.onChangeCallCount)
				c.input.AddPosting()
				assert.Equal(t, 2, c.input.CountPostings())
				assert.Equal(t, 2, c.onChangeCallCount)
				// Advance one and delete it
				c.input.DeleteCurrentPosting()
				assert.Equal(t, 1, c.input.CountPostings())
				assert.Equal(t, 3, c.onChangeCallCount)
				// Delete last one
				c.input.DeleteCurrentPosting()
				assert.Equal(t, 0, c.input.CountPostings())
				assert.Equal(t, 4, c.onChangeCallCount)
				// Last delete does nothing
				c.input.DeleteCurrentPosting()
				assert.Equal(t, 0, c.input.CountPostings())
				assert.Equal(t, 4, c.onChangeCallCount)
			},
		},
		{
			"Current Posting",
			func(t *testing.T, c *context) {
				addedPosting := c.input.CurrentPosting()
				assert.Equal(t, 2, c.onChangeCallCount)
				currentPosting := c.input.CurrentPosting()
				assert.Same(t, addedPosting, currentPosting)
			},
		},
		{
			"Calculate posting balance no postings",
			func(t *testing.T, c *context) {
				expected := []journal.Ammount{}
				assert.Equal(t, expected, c.input.PostingBalance())
			},
		},
		{
			"Calculate posting balance with postings total 0",
			func(t *testing.T, c *context) {
				ammount1 := journal.Ammount{
					Commodity: "EUR",
					Quantity:  decimal.New(12, 1),
				}
				c.input.AddPosting().SetAmmount(ammount1)
				ammount2 := journal.Ammount{
					Commodity: "EUR",
					Quantity:  decimal.New(-12, 1),
				}
				c.input.AddPosting().SetAmmount(ammount2)
				expected := []journal.Ammount{}
				assert.Equal(t, expected, c.input.PostingBalance())
			},
		},
		{
			"Calculate posting balance with postings total not 0",
			func(t *testing.T, c *context) {
				ammount1 := journal.Ammount{
					Commodity: "EUR",
					Quantity:  decimal.New(12, 1),
				}
				c.input.AddPosting().SetAmmount(ammount1)
				ammount2 := journal.Ammount{
					Commodity: "BRL",
					Quantity:  decimal.New(-12, 1),
				}
				c.input.AddPosting().SetAmmount(ammount2)
				expected := []journal.Ammount{ammount1, ammount2}
				assert.ElementsMatch(t, expected, c.input.PostingBalance())
			},
		},
		{
			"Ignore postings without ammount",
			func(t *testing.T, c *context) {
				c.input.AddPosting()
				c.input.AddPosting()
				c.input.AddPosting().SetAmmount(anAmmount)
				expected := []journal.Ammount{anAmmount}
				assert.Equal(t, expected, c.input.PostingBalance())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := NewJournalEntryInput()
			ctx := &context{false, 0, input}
			input.AddOnChangeHook(func() { ctx.onChangeCalled = true })
			input.AddOnChangeHook(func() { ctx.onChangeCallCount++ })
			tc.run(t, ctx)
		})
	}
}

func TestRepr(t *testing.T) {
	type testcase struct {
		name         string
		journalEntry func(t *testing.T) *JournalEntryInput
		expected     string
	}
	testcases := []testcase{
		{
			name: "Empty",
			journalEntry: func(_ *testing.T) *JournalEntryInput {
				return &JournalEntryInput{}
			},
			expected: "",
		},
		{
			name: "with date",
			journalEntry: func(t *testing.T) *JournalEntryInput {
				i := NewJournalEntryInput()
				i.SetDate(tu.Date1(t))
				return i
			},
			expected: "1993-11-23",
		},
		{
			name: "with description",
			journalEntry: func(t *testing.T) *JournalEntryInput {
				i := NewJournalEntryInput()
				i.SetDate(tu.Date1(t))
				i.SetDescription("FOO BAR")
				return i
			},
			expected: "1993-11-23 FOO BAR",
		},
		{
			name: "with postings",
			journalEntry: func(t *testing.T) *JournalEntryInput {
				i := NewJournalEntryInput()
				i.SetDate(tu.Date1(t))
				i.SetDescription("FOO BAR")
				posting1 := i.AddPosting()
				posting1.SetAccount("ACC")
				posting1.SetAmmount(anAmmount)
				posting2 := i.AddPosting()
				posting2.SetAccount("ACC2")
				posting2.SetAmmount(anotherAmmount)
				return i
			},
			expected: strings.Join(
				[]string{
					"1993-11-23 FOO BAR",
					"    ACC    EUR 2.2",
					"    ACC2    EUR -2.2",
				},
				"\n",
			),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.journalEntry(t).Repr()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTextToAmmount(t *testing.T) {
	type testcase struct {
		text     string
		ammount  journal.Ammount
		errorMsg string
	}
	var testcases = []testcase{
		{
			text: "EUR 12.20",
			ammount: journal.Ammount{
				Commodity: "EUR",
				Quantity:  decimal.New(1220, -2),
			},
		},
		{
			text: "EUR 99999.99999",
			ammount: journal.Ammount{
				Commodity: "EUR",
				Quantity:  decimal.NewFromFloat(99999.99999),
			},
		},
		{
			text: "12.20",
			ammount: journal.Ammount{
				Commodity: "",
				Quantity:  decimal.New(1220, -2),
			},
		},
		{
			text:     "12,20",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR 12 12",
			errorMsg: "invalid format",
		},
		{
			text:     "12 FOO",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR  12.20",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR 12.20 ",
			errorMsg: "invalid format",
		},
		{
			text:     " EUR 12.20 ",
			errorMsg: "invalid format",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.text, func(t *testing.T) {
			result, err := TextToAmmount(tc.text)
			if tc.errorMsg == "" {
				assert.Nil(t, err)
				assert.Equal(t, tc.ammount, result)
			} else {
				assert.ErrorContains(t, err, tc.errorMsg)
			}
		})
	}
}

func TestTextToPosting(t *testing.T) {
	type testcase struct {
		name                string
		input               func(t *testing.T) *JournalEntryInput
		expectedTransaction func(t *testing.T) *journal.Transaction
		errorMsg            string
	}
	var testcases = []testcase{
		{
			name: "Simple transaction",
			input: func(t *testing.T) *JournalEntryInput {
				return tu.JournalEntryInput_1(t)
			},
			expectedTransaction: func(t *testing.T) *journal.Transaction {
				return tu.Transaction_1(t)
			},
		},
		{
			name: "Missing description",
			input: func(t *testing.T) *JournalEntryInput {
				out := tu.JournalEntryInput_1(t)
				out.ClearDescription()
				return out
			},
			errorMsg: "missing description",
		},
		{
			name: "Missing date",
			input: func(t *testing.T) *JournalEntryInput {
				out := tu.JournalEntryInput_1(t)
				out.ClearDate()
				return out
			},
			errorMsg: "missing date",
		},
		{
			name: "Unbalanced posting",
			input: func(t *testing.T) *JournalEntryInput {
				out := tu.JournalEntryInput_1(t)
				posting_3 := out.AddPosting()
				tu.FillPostingInput_3(t, posting_3)
				return out
			},
			errorMsg: "postings are not balanced",
		},
		{
			name: "Posting missing ammount",
			input: func(t *testing.T) *JournalEntryInput {
				out := tu.JournalEntryInput_1(t)
				posting := out.AddPosting()
				posting.SetAccount("FOO")
				return out
			},
			errorMsg: "one of the postings is missing the amount",
		},
		{
			name: "Posting missing account",
			input: func(t *testing.T) *JournalEntryInput {
				out := tu.JournalEntryInput_1(t)
				posting, _ := out.GetPosting(0)
				posting.ClearAccount()
				return out
			},
			errorMsg: "one of the postings is missing the account",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			input := tc.input(t)
			result, err := input.ToTransaction()
			if tc.errorMsg == "" {
				assert.Nil(t, err)
				assert.Equal(t, *tc.expectedTransaction(t), result)
			} else {
				assert.ErrorContains(t, err, tc.errorMsg)
			}
		})
	}
}
