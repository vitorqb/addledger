package accountguesser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/accountguesser"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/testutils"
)

func TestAccountGuesser(t *testing.T) {
	type testcontext struct {
		accountguesser *AccountGuesser
	}
	type testcase struct {
		name               string
		transactionHistory func() TransactionHistory
		inputPostings      func() []journal.Posting
		description        string
		success            bool
		expected           journal.Account
	}
	var testcases = []testcase{
		{
			name: "no transaction history",
			transactionHistory: func() TransactionHistory {
				return []journal.Transaction{}
			},
			inputPostings: func() []journal.Posting { return []journal.Posting{} },
			description:   "Supermarket",
			success:       false,
		},
		{
			name: "perfect match",
			transactionHistory: func() TransactionHistory {
				return []journal.Transaction{
					{
						Description: "Supermarket",
						Date:        testutils.Date1(t),
						Posting: []journal.Posting{
							{
								Account: "expenses:supermarket",
								Ammount: journal.Ammount{},
							},
						},
					},
				}
			},
			inputPostings: func() []journal.Posting { return []journal.Posting{} },
			description:   "Supermarket",
			success:       true,
			expected:      "expenses:supermarket",
		},
		{
			name: "close match",
			transactionHistory: func() TransactionHistory {
				return []journal.Transaction{
					{
						Description: "Supermarkaa",
						Date:        testutils.Date1(t),
						Posting: []journal.Posting{
							{
								Account: "expenses:supermarket",
								Ammount: journal.Ammount{},
							},
						},
					},
				}
			},
			inputPostings: func() []journal.Posting {
				return []journal.Posting{}
			},
			description: "Supermarket",
			success:     true,
			expected:    "expenses:supermarket",
		},
		{
			name: "with previous input posting",
			transactionHistory: func() TransactionHistory {
				return []journal.Transaction{
					{
						Description: "Supermarket",
						Date:        testutils.Date1(t),
						Posting: []journal.Posting{
							{
								Account: "expenses:supermarket",
								Ammount: journal.Ammount{},
							},
							{
								Account: "assets:current-account",
								Ammount: journal.Ammount{},
							},
						},
					},
				}
			},
			inputPostings: func() []journal.Posting {
				return []journal.Posting{
					{
						Account: "expenses:supermarket",
						Ammount: journal.Ammount{},
					},
				}
			},
			description: "Supermarket",
			success:     true,
			expected:    "assets:current-account",
		},
		{
			name: "gets perfect match over close match",
			transactionHistory: func() TransactionHistory {
				return []journal.Transaction{
					{
						Description: "Supermarke",
						Date:        testutils.Date1(t),
						Posting: []journal.Posting{
							{
								Account: "foo",
								Ammount: journal.Ammount{},
							},
							{
								Account: "bar",
								Ammount: journal.Ammount{},
							},
						},
					},
					{
						Description: "Supermarket",
						Date:        testutils.Date1(t),
						Posting: []journal.Posting{
							{
								Account: "expenses:supermarket",
								Ammount: journal.Ammount{},
							},
							{
								Account: "assets:current-account",
								Ammount: journal.Ammount{},
							},
						},
					},
				}
			},
			inputPostings: func() []journal.Posting {
				return []journal.Posting{}
			},
			description: "Supermarket",
			success:     true,
			expected:    "expenses:supermarket",
		},
		{
			name: "if two matches get most recent one",
			transactionHistory: func() TransactionHistory {
				oldTransaction := journal.Transaction{
					Description: "Supermarket",
					Date:        testutils.Date1(t),
					Posting: []journal.Posting{
						{
							Account: "foo",
							Ammount: journal.Ammount{},
						},
						{
							Account: "bar",
							Ammount: journal.Ammount{},
						},
					},
				}
				recentTransaction := oldTransaction
				recentTransaction.Date = testutils.Date2(t)
				recentTransaction.Posting[0].Account = "expenses:supermarket"
				return []journal.Transaction{oldTransaction, recentTransaction}
			},
			inputPostings: func() []journal.Posting {
				return []journal.Posting{}
			},
			description: "Supermarket",
			success:     true,
			expected:    "expenses:supermarket",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := new(testcontext)
			c.accountguesser = New()
			actual, success := c.accountguesser.Guess(tc.transactionHistory(), tc.inputPostings(), tc.description)
			assert.Equal(t, tc.success, success)
			if tc.success {
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}
