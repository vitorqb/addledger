package state_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
	. "github.com/vitorqb/addledger/internal/state"
)

func TestPostingData(t *testing.T) {
	type testcontext struct {
		data              *PostingData
		onChangeCallCount int
	}
	type testcase struct {
		name string
		run  func(*testing.T, *testcontext)
	}
	testcases := []testcase{
		{
			name: "Notifies when account changes",
			run: func(t *testing.T, ctx *testcontext) {
				ctx.data.Account.Set(journal.Account("foo"))
				assert.Equal(t, 1, ctx.onChangeCallCount)
				ctx.data.Account.Clear()
				assert.Equal(t, 2, ctx.onChangeCallCount)
			},
		},
		{
			name: "Notifies when ammount changes",
			run: func(t *testing.T, ctx *testcontext) {
				ctx.data.Ammount.Set(finance.Ammount{})
				assert.Equal(t, 1, ctx.onChangeCallCount)
				ctx.data.Ammount.Clear()
				assert.Equal(t, 2, ctx.onChangeCallCount)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &testcontext{}
			ctx.data = NewPostingData()
			ctx.data.AddOnChangeHook(func() { ctx.onChangeCallCount++ })
			tc.run(t, ctx)
		})
	}
}

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
		{
			name: "Notifies when tags change",
			run: func(t *testing.T, ctx *testcontext) {
				tag := journal.Tag{Name: "foo", Value: "bar"}
				ctx.data.Tags.Append(tag)
				assert.Equal(t, 1, ctx.onChangeCallCount)
				ctx.data.Tags.Pop()
				assert.Equal(t, 2, ctx.onChangeCallCount)
			},
		},
		{
			name: "Notifies when postings change",
			run: func(t *testing.T, ctx *testcontext) {
				ctx.data.Postings.Append(NewPostingData())
				assert.Equal(t, 1, ctx.onChangeCallCount)
				ctx.data.Postings.Pop()
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
