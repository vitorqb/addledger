package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostingInput(t *testing.T) {
	t.Run("Set Account", func(t *testing.T) {
		// !!!! TODO Make it DRY
		onChangeCallCount := 0
		onChange := func() { onChangeCallCount++ }
		postingInput := NewPostingInput()
		postingInput.AddOnChangeHook(onChange)
		account, found := postingInput.GetAccount()
		assert.False(t, found)
		postingInput.SetAccount("FOO")
		account, found = postingInput.GetAccount()
		assert.True(t, found)
		assert.Equal(t, "FOO", account)
		assert.Equal(t, 1, onChangeCallCount)
	})
}
