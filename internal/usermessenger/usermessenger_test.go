package usermessenger_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	statemod "github.com/vitorqb/addledger/internal/state"
	. "github.com/vitorqb/addledger/internal/usermessenger"
)

func TestUserMessenger(t *testing.T) {
	t.Run("Info", func(t *testing.T) {
		state := statemod.InitialState()
		uMessenger := New(state)
		uMessenger.Info("info")
		assert.Equal(t, "info", state.Display.UserMessage())
	})
	t.Run("Warning", func(t *testing.T) {
		t.Run("no error", func(t *testing.T) {
			state := statemod.InitialState()
			uMessenger := New(state)
			uMessenger.Warning("warning", nil)
			assert.Equal(t, "WARNING: warning", state.Display.UserMessage())
		})
		t.Run("with error", func(t *testing.T) {
			state := statemod.InitialState()
			uMessenger := New(state)
			uMessenger.Warning("warning", fmt.Errorf("error"))
			assert.Equal(t, "WARNING: warning - error", state.Display.UserMessage())
		})
	})
	t.Run("Error", func(t *testing.T) {
		t.Run("no error", func(t *testing.T) {
			state := statemod.InitialState()
			uMessenger := New(state)
			uMessenger.Error("error", nil)
			assert.Equal(t, "ERROR: error", state.Display.UserMessage())
		})
		t.Run("with error", func(t *testing.T) {
			state := statemod.InitialState()
			uMessenger := New(state)
			uMessenger.Error("error", fmt.Errorf("error"))
			assert.Equal(t, "ERROR: error - error", state.Display.UserMessage())
		})
	})

}
