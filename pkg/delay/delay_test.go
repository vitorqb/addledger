package delay_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/pkg/delay"
)

func TestFunction(t *testing.T) {

	t.Run("Test running with delay", func(t *testing.T) {
		calls := 0
		fun := func() { calls++ }
		run := NewFunction(fun, 250*time.Millisecond)
		tick := run.Tick()
		run.Schedule()
		run.Schedule()
		<-tick
		assert.Equal(t, 1, calls)
	})

}
