package userinput_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/testutils"
	. "github.com/vitorqb/addledger/internal/userinput"
)

func TestTransactionRepr(t *testing.T) {
	type testcase struct {
		name        string
		transaction func(*testing.T, *state.TransactionData)
		expected    string
	}
	testcases := []testcase{
		{
			name:        "Empty transaction",
			transaction: func(_ *testing.T, tra *state.TransactionData) {},
			expected:    "",
		},
		{
			name: "With date",
			transaction: func(_ *testing.T, tra *state.TransactionData) {
				tra.Date.Set(testutils.Date1(t))
			},
			expected: "1993-11-23",
		},
		{
			name: "With description",
			transaction: func(_ *testing.T, tra *state.TransactionData) {
				tra.Date.Set(testutils.Date1(t))
				tra.Description.Set("foo")
			},
			expected: "1993-11-23 foo",
		},
		{
			name: "With tags",
			transaction: func(_ *testing.T, tra *state.TransactionData) {
				tra.Date.Set(testutils.Date1(t))
				tra.Description.Set("foo")
				tra.Tags.Append(journal.Tag{Name: "bar", Value: "baz"})
			},
			expected: "1993-11-23 foo ; bar:baz",
		},
		{
			name: "With postings",
			transaction: func(_ *testing.T, tra *state.TransactionData) {
				tra.Date.Set(testutils.Date1(t))
				tra.Description.Set("foo")
				tra.Tags.Append(journal.Tag{Name: "bar", Value: "baz"})
				posting := state.NewPostingData()
				posting.Account.Set(journal.Account("ACC"))
				posting.Ammount.Set(*testutils.Ammount_1(t))
				tra.Postings.Append(posting)
				posting2 := state.NewPostingData()
				posting2.Account.Set(journal.Account("ACC2"))
				posting2.Ammount.Set((*testutils.Ammount_1(t)).InvertSign())
				tra.Postings.Append(posting2)
			},
			expected: strings.Join([]string{
				"1993-11-23 foo ; bar:baz",
				"    ACC    EUR 2.2",
				"    ACC2    EUR -2.2",
			}, "\n"),
		},
		{
			name: "With posting without commodity",
			transaction: func(_ *testing.T, tra *state.TransactionData) {
				tra.Date.Set(testutils.Date1(t))
				posting := state.NewPostingData()
				amount := testutils.Ammount_1(t)
				amount.Commodity = ""
				posting.Account.Set(journal.Account("ACC"))
				posting.Ammount.Set(*amount)
				tra.Postings.Append(posting)
			},
			expected: strings.Join([]string{
				"1993-11-23",
				"    ACC    2.2",
			}, "\n"),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			trans := state.NewTransactionData()
			tc.transaction(t, trans)
			actual := TransactionRepr(trans)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
