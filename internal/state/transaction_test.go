package state_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/state"
)

func TestTransactionData(t *testing.T) {
	type testcontext struct {
		data              *TransactionData
		onChangeCallCount int
	}
	type testcase struct {
		name string
		run  func(*testing.T, *testcontext)
	}
	testcases := []testcase{
		{
			name: "Notifies when date changes",
			run: func(t *testing.T, ctx *testcontext) {
				ctx.data.Date.Set(time.Now())
				assert.Equal(t, 1, ctx.onChangeCallCount)
				ctx.data.Date.Clear()
				assert.Equal(t, 2, ctx.onChangeCallCount)
			},
		},
		{
			name: "Notifies when description changes",
			run: func(t *testing.T, ctx *testcontext) {
				ctx.data.Description.Set("foo")
				assert.Equal(t, 1, ctx.onChangeCallCount)
				ctx.data.Description.Clear()
				assert.Equal(t, 2, ctx.onChangeCallCount)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &testcontext{}
			ctx.data = NewTransactionData()
			ctx.data.AddOnChangeHook(func() { ctx.onChangeCallCount++ })
			tc.run(t, ctx)
		})
	}
}
