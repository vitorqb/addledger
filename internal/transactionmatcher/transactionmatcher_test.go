package transactionmatcher_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/stringmatcher"
	. "github.com/vitorqb/addledger/internal/transactionmatcher"
)

func TestTransactionMatcher(t *testing.T) {

	type testcontext struct {
		transactionMatcher *TransactionMatcher
	}

	type testcase struct {
		name string
		run  func(*testing.T, *testcontext)
	}

	testcases := []testcase{
		{
			name: "Simple Match",
			run: func(t *testing.T, ctx *testcontext) {
				ctx.transactionMatcher.SetDescriptionInput("test")
				transactions := []journal.Transaction{{Description: "test"}}
				ctx.transactionMatcher.SetTransactionHistory(transactions)
				matches := ctx.transactionMatcher.Match()
				assert.Equal(t, transactions, matches)
			},
		},
		{
			name: "No match _ no history",
			run: func(t *testing.T, ctx *testcontext) {
				ctx.transactionMatcher.SetDescriptionInput("test")
				ctx.transactionMatcher.SetTransactionHistory([]journal.Transaction{})
				matches := ctx.transactionMatcher.Match()
				assert.Equal(t, []journal.Transaction{}, matches)
			},
		},
		{
			name: "No match _ description doesnt match",
			run: func(t *testing.T, ctx *testcontext) {
				ctx.transactionMatcher.SetDescriptionInput("test")
				transactions := []journal.Transaction{{Description: "i dont match"}}
				ctx.transactionMatcher.SetTransactionHistory(transactions)
				matches := ctx.transactionMatcher.Match()
				assert.Equal(t, []journal.Transaction{}, matches)
			},
		},
		{
			name: "Order matches by distance",
			run: func(t *testing.T, ctx *testcontext) {
				ctx.transactionMatcher.SetDescriptionInput("test123")
				transactions := []journal.Transaction{
					{Description: "test100"}, // distance 2
					{Description: "test120"}, // distance 1
					{Description: "test123"}, // distance 0
					{Description: "test000"}, // distance 3
				}
				sortedTransactions := []journal.Transaction{
					{Description: "test123"}, // distance 0
					{Description: "test120"}, // distance 1
					{Description: "test100"}, // distance 2
					{Description: "test000"}, // distance 3
				}
				ctx.transactionMatcher.SetTransactionHistory(transactions)
				matches := ctx.transactionMatcher.Match()
				assert.Equal(t, sortedTransactions, matches)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Note: we could use a mock here, but I think it's not worth it.
			stringMatcher, err := stringmatcher.New(&stringmatcher.Options{})
			if err != nil {
				t.Fatal(err)
			}
			transactionMatcher := New(stringMatcher)
			ctx := &testcontext{transactionMatcher: transactionMatcher}
			tc.run(t, ctx)
		})
	}
}
