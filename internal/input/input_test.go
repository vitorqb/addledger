package input

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert.True(t, c.onChangeCalled)
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
				assert.True(t, c.onChangeCalled)
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
			"Current Posting",
			func(t *testing.T, c *context) {
				addedPosting := c.input.CurrentPosting()
				assert.Equal(t, 2, c.onChangeCallCount)
				currentPosting := c.input.CurrentPosting()
				assert.Same(t, addedPosting, currentPosting)
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
				posting1.SetValue("EUR 12.20")
				posting2 := i.AddPosting()
				posting2.SetAccount("ACC2")
				posting2.SetValue("EUR -12.20")
				return i
			},
			expected: strings.Join(
				[]string{
					"1993-11-23 FOO BAR",
					"    ACC    EUR 12.20",
					"    ACC2    EUR -12.20",
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
