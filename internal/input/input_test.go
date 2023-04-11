package input

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var aDate, _ = time.Parse("2006-01-02", "1993-11-23")

func TestJournalEntryInput(t *testing.T) {

	type context struct {
		onChangeCalled bool
		onChangeCallCount int
		input *JournalEntryInput
	}

	type test struct {
		name string
		run func(*testing.T, *context)
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
				c.input.AddPosting()
				posting, found := c.input.GetPosting(0)
				assert.True(t, found)
				_, found = posting.GetAccount()
				assert.False(t, found)
				posting.SetAccount("FOO")
				account, found := posting.GetAccount()
				assert.True(t, found)
				assert.Equal(t, account, "FOO")
				assert.Equal(t, 2, c.onChangeCallCount)
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
