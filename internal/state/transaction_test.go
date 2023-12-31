package state_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/state"
)

func TestTransactionData(t *testing.T) {
	t.Run("Notifies when date changes", func(t *testing.T) {
		onChangeCallCount := 0
		data := NewTransactionData()
		data.AddOnChangeHook(func() { onChangeCallCount++ })
		data.Date.Set(time.Now())
		assert.Equal(t, 1, onChangeCallCount)
		data.Date.Clear()
		assert.Equal(t, 2, onChangeCallCount)
	})
}
