package userinput_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/testutils"
	tu "github.com/vitorqb/addledger/internal/testutils"
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

func TestTransactionFromData(t *testing.T) {
	type testcase struct {
		name                string
		data                func(t *testing.T) *state.TransactionData
		expectedTransaction func(t *testing.T) *journal.Transaction
		errorMsg            string
	}
	var testcases = []testcase{
		{
			name: "Simple transaction",
			data: func(t *testing.T) *state.TransactionData {
				return tu.TransactionData_1(t)
			},
			expectedTransaction: func(t *testing.T) *journal.Transaction {
				return tu.Transaction_1(t)
			},
		},
		{
			name: "Simple transaction with tags",
			data: func(t *testing.T) *state.TransactionData {
				out := tu.TransactionData_1(t)
				tag := journal.Tag{Name: "foo", Value: "bar"}
				out.Tags.Append(tag)
				return out
			},
			expectedTransaction: func(t *testing.T) *journal.Transaction {
				out := tu.Transaction_1(t)
				out.Comment = "foo:bar"
				return out
			},
		},
		{
			name: "Missing description",
			data: func(t *testing.T) *state.TransactionData {
				out := tu.TransactionData_1(t)
				out.Description.Clear()
				return out
			},
			errorMsg: "missing description",
		},
		{
			name: "Missing date",
			data: func(t *testing.T) *state.TransactionData {
				out := tu.TransactionData_1(t)
				out.Date.Clear()
				return out
			},
			errorMsg: "missing date",
		},
		{
			name: "Unbalanced posting",
			data: func(t *testing.T) *state.TransactionData {
				out := tu.TransactionData_1(t)
				posting_3 := state.NewPostingData()
				posting_3.Account.Set("ACC1")
				posting_3.Ammount.Set(*testutils.Ammount_1(t))
				out.Postings.Append(posting_3)
				return out
			},
			errorMsg: "postings are not balanced",
		},
		{
			name: "Posting missing ammount",
			data: func(t *testing.T) *state.TransactionData {
				out := tu.TransactionData_1(t)
				posting, _ := out.Postings.Last()
				posting.Ammount.Clear()
				return out
			},
			errorMsg: "one of the postings is missing the ammount",
		},
		{
			name: "Posting missing account",
			data: func(t *testing.T) *state.TransactionData {
				out := tu.TransactionData_1(t)
				posting, _ := out.Postings.Last()
				posting.Account.Clear()
				return out
			},
			errorMsg: "one of the postings is missing the account",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			data := tc.data(t)
			result, err := TransactionFromData(data)
			if tc.errorMsg == "" {
				assert.Nil(t, err)
				assert.Equal(t, *tc.expectedTransaction(t), result)
			} else {
				assert.ErrorContains(t, err, tc.errorMsg)
			}
		})
	}
}

func TestRemoveIncompletePostings(t *testing.T) {
	type testcase struct {
		name     string
		setup    func(t *testing.T) []*state.PostingData
		expected []journal.Posting
	}
	var testcases = []testcase{
		{
			name: "No postings",
			setup: func(t *testing.T) []*state.PostingData {
				return []*state.PostingData{}
			},
		},
		{
			name: "One complete posting",
			setup: func(t *testing.T) []*state.PostingData {
				posting := state.NewPostingData()
				posting.Account.Set("ACC")
				posting.Ammount.Set(*testutils.Ammount_1(t))
				return []*state.PostingData{posting}
			},
			expected: []journal.Posting{
				{
					Account: "ACC",
					Ammount: *testutils.Ammount_1(t),
				},
			},
		},
		{
			name: "One missing ammount",
			setup: func(t *testing.T) []*state.PostingData {
				posting := state.NewPostingData()
				posting.Account.Set("ACC")
				return []*state.PostingData{posting}
			},
		},
		{
			name: "One missing account",
			setup: func(t *testing.T) []*state.PostingData {
				posting := state.NewPostingData()
				posting.Ammount.Set(*testutils.Ammount_1(t))
				return []*state.PostingData{posting}
			},
		},
		{
			name: "One missing account and one complete",
			setup: func(t *testing.T) []*state.PostingData {
				posting_1 := state.NewPostingData()
				posting_1.Ammount.Set(*testutils.Ammount_1(t))
				posting_2 := state.NewPostingData()
				posting_2.Account.Set("ACC")
				posting_2.Ammount.Set(*testutils.Ammount_1(t))
				return []*state.PostingData{posting_1, posting_2}
			},
			expected: []journal.Posting{
				{
					Account: "ACC",
					Ammount: *testutils.Ammount_1(t),
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			postings := tc.setup(t)
			result := ExtractPostings(postings)
			assert.Equal(t, tc.expected, result)
		})
	}
}
