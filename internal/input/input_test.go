package input

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var aDate, _ = time.Parse("2006-01-02", "1993-11-23")

func TestJournalEntryInput(t *testing.T) {

	t.Run("Date", func(t *testing.T) {
		onChangeCalled := false
		onChange := func() { onChangeCalled = true }
		input := NewJournalEntryInput()
		input.AddOnChangeHook(onChange)
		_, found := input.GetDate()
		assert.False(t, found)
		input.SetDate(aDate)
		date, found := input.GetDate()
		assert.True(t, found)
		assert.Equal(t, date, aDate)
		assert.True(t, onChangeCalled)
	})

	t.Run("Description", func(t *testing.T) {
		onChangeCalled := false
		onChange := func() { onChangeCalled = true }
		input := NewJournalEntryInput()
		input.AddOnChangeHook(onChange)
		_, found := input.GetDescription()
		assert.False(t, found)
		input.SetDescription("FOO")
		description, found := input.GetDescription()
		assert.True(t, found)
		assert.Equal(t, description, "FOO")
		assert.True(t, onChangeCalled)
	})
}
