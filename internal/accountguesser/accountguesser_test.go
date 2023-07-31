package accountguesser_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/accountguesser"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/stringmatcher"
	"github.com/vitorqb/addledger/internal/testutils"
	. "github.com/vitorqb/addledger/mocks/stringmatcher"
)

func TestAccountGuesser(t *testing.T) {
	type testcontext struct {
		accountguesser *AccountGuesser
		stringMatcher  *MockIStringMatcher
	}
	type testcase struct {
		name               string
		setup              func(*testing.T, *testcontext)
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
			name: "distance higher than 15 (from cache)",
			transactionHistory: func() TransactionHistory {
				return []journal.Transaction{
					{
						Description: "BA",
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
			setup: func(t1 *testing.T, t2 *testcontext) {
				t2.stringMatcher.EXPECT().Distance("AB", "BA").Return(999)
			},
			inputPostings: func() []journal.Posting { return []journal.Posting{} },
			description:   "AB",
			success:       false,
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

	// Default setup if test does not define one
	defaultSetup := func(t *testing.T, c *testcontext) {
		// Ensure stringMatcher.Distance return real distance
		realMacher, err := stringmatcher.New(&stringmatcher.Options{})
		if err != nil {
			t.Fatal(err)
		}
		c.stringMatcher.EXPECT().Distance(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(a string, b string) int {
			return realMacher.Distance(a, b)
		})
	}

	// Run test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.stringMatcher = NewMockIStringMatcher(ctrl)
			if tc.setup != nil {
				tc.setup(t, c)
			} else {
				defaultSetup(t, c)
			}
			c.accountguesser, err = New(Options{StringMatcher: c.stringMatcher})
			if err != nil {
				t.Fatal(err)
			}
			actual, success := c.accountguesser.Guess(tc.transactionHistory(), tc.inputPostings(), tc.description)
			assert.Equal(t, tc.success, success)
			if tc.success {
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}
