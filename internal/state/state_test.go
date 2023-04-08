package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {

	t.Run("Notify on change of JournalEntryInput", func(t *testing.T) {
		hookCalled := false
		hook := func() { hookCalled = true }
		s := InitialState()
		s.AddOnChangeHook(hook)
		s.JournalEntryInput.SetDescription("FOO")
		assert.True(t, hookCalled)
	})

	t.Run("NextPhase", func(t *testing.T) {
		hookCalled := false
		hook := func() { hookCalled = true }
		s := InitialState()
		s.AddOnChangeHook(hook)
		assert.Equal(t, s.CurrentPhase, InputDate)
		s.NextPhase()
		assert.Equal(t, s.CurrentPhase, InputDescription)
		assert.True(t, hookCalled)
	})
}
