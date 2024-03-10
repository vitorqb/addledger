package userinput_test

import (
	"strings"
	"testing"

	"fmt"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/finance"
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

func TestTagTagToText(t *testing.T) {
	tag := journal.Tag{
		Name:  "foo",
		Value: "bar",
	}
	expected := "foo:bar"
	actual := TagToText(tag)
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	tags := []journal.Tag{
		{Name: "foo", Value: "bar"},
		{Name: "baz", Value: "qux"},
	}
	expectedArr := []string{"foo:bar", "baz:qux"}
	actualArr := TagsToText(tags)
	if actual != expected {
		t.Errorf("Expected %s, got %s", expectedArr, actualArr)
	}
}

func TestTagTextToTag__Good(t *testing.T) {
	type testCase struct {
		input    string
		expected journal.Tag
	}
	for _, tc := range []testCase{
		{
			input:    "foo:bar",
			expected: journal.Tag{Name: "foo", Value: "bar"},
		},
		{
			input:    " foo:bar ",
			expected: journal.Tag{Name: "foo", Value: "bar"},
		},
		{
			input:    "foo-bar:baz",
			expected: journal.Tag{Name: "foo-bar", Value: "baz"},
		},
		{
			input:    "foo_bar:baz",
			expected: journal.Tag{Name: "foo_bar", Value: "baz"},
		},
		{
			input:    "foo_bar:baz_123",
			expected: journal.Tag{Name: "foo_bar", Value: "baz_123"},
		},
		{
			input:    "foo_bar:baz-123",
			expected: journal.Tag{Name: "foo_bar", Value: "baz-123"},
		},
		{
			input:    "foo_bar:baz-123_abc",
			expected: journal.Tag{Name: "foo_bar", Value: "baz-123_abc"},
		},
	} {
		tag, err := TextToTag(tc.input)
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}
		if tag != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, tag)
		}
	}
}

func TestTagTextToTag__Bad(t *testing.T) {
	for _, input := range []string{
		"foo",
		"foo:",
		"foo:bar:baz",
		"",
		"foo bar:baz",
		"foo:bar baz",
		"some word",
		"some word:foo",
	} {
		_, err := TextToTag(input)
		if err == nil {
			t.Errorf(fmt.Sprintf("Expected error, got none: %s", input))
		}

	}
}

func TestTextToAmmount(t *testing.T) {
	type testcase struct {
		text     string
		ammount  finance.Ammount
		errorMsg string
	}
	var testcases = []testcase{
		{
			text: "EUR 12.20",
			ammount: finance.Ammount{
				Commodity: "EUR",
				Quantity:  decimal.New(1220, -2),
			},
		},
		{
			text: "EUR 99999.99999",
			ammount: finance.Ammount{
				Commodity: "EUR",
				Quantity:  decimal.NewFromFloat(99999.99999),
			},
		},
		{
			text: "12.20",
			ammount: finance.Ammount{
				Commodity: "",
				Quantity:  decimal.New(1220, -2),
			},
		},
		{
			text:     "12,20",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR 12 12",
			errorMsg: "invalid format",
		},
		{
			text:     "12 FOO",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR  12.20",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR 12.20 ",
			errorMsg: "invalid format",
		},
		{
			text:     " EUR 12.20 ",
			errorMsg: "invalid format",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.text, func(t *testing.T) {
			result, err := TextToAmmount(tc.text)
			if tc.errorMsg == "" {
				assert.Nil(t, err)
				assert.Equal(t, tc.ammount, result)
			} else {
				assert.ErrorContains(t, err, tc.errorMsg)
			}
		})
	}
}

func TestPostingFromData(t *testing.T) {
	t.Run("Missing ammount", func(t *testing.T) {
		data := state.NewPostingData()
		_, err := PostingFromData(data)
		assert.ErrorContains(t, err, "missing the ammount")
	})
	t.Run("Missing account", func(t *testing.T) {
		data := state.NewPostingData()
		data.Ammount.Set(*testutils.Ammount_1(t))
		_, err := PostingFromData(data)
		assert.ErrorContains(t, err, "missing the account")
	})
	t.Run("Complete", func(t *testing.T) {
		ammount := testutils.Ammount_1(t)
		data := state.NewPostingData()
		data.Ammount.Set(*ammount)
		data.Account.Set("ACC")
		posting, err := PostingFromData(data)
		assert.Nil(t, err)
		expPosting := journal.Posting{Account: "ACC", Ammount: *ammount}
		assert.Equal(t, expPosting, posting)
	})
}

func TestPostingsFromData(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		data := []*state.PostingData{}
		postings, err := PostingsFromData(data)
		assert.Nil(t, err)
		assert.Equal(t, []journal.Posting{}, postings)
	})
	t.Run("Two complete", func(t *testing.T) {
		ammount_1 := testutils.Ammount_1(t)
		ammount_2 := testutils.Ammount_2(t)
		data := make([]*state.PostingData, 2)
		data[0] = state.NewPostingData()
		data[0].Ammount.Set(*ammount_1)
		data[0].Account.Set("ACC1")
		data[1] = state.NewPostingData()
		data[1].Ammount.Set(*ammount_2)
		data[1].Account.Set("ACC2")
		postings, err := PostingsFromData(data)
		assert.Nil(t, err)
		expPostings := []journal.Posting{
			{Ammount: *ammount_1, Account: "ACC1"},
			{Ammount: *ammount_2, Account: "ACC2"},
		}
		assert.Equal(t, expPostings, postings)
	})
	t.Run("One missing ammount", func(t *testing.T) {
		data := make([]*state.PostingData, 1)
		data[0] = state.NewPostingData()
		_, err := PostingsFromData(data)
		assert.ErrorContains(t, err, "missing the ammount")
	})
	t.Run("One missing account", func(t *testing.T) {
		ammount := testutils.Ammount_1(t)
		data := make([]*state.PostingData, 1)
		data[0] = state.NewPostingData()
		data[0].Ammount.Set(*ammount)
		_, err := PostingsFromData(data)
		assert.ErrorContains(t, err, "missing the account")
	})
}
