package input

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var aDate, _ = time.Parse("2006-01-02", "1993-11-23")

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
				c.input.SetDate(aDate)
				date, found := c.input.GetDate()
				assert.True(t, found)
				assert.Equal(t, date, aDate)
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
