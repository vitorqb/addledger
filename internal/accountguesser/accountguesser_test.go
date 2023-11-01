package accountguesser_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/accountguesser"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
	tu "github.com/vitorqb/addledger/internal/testutils"
	. "github.com/vitorqb/addledger/mocks/accountguesser"
)

func TestMatchedTransactionsGuesser(t *testing.T) {
	type testcontext struct {
		accountguesser *MatchedTransactionsGuesser
	}
	type testcase struct {
		name                string
		setup               func(*testing.T, *testcontext)
		matchedTransactions func(*testing.T) MatchedTransactions
		inputPostings       func() []journal.Posting
		success             bool
		expected            journal.Account
	}
	var testcases = []testcase{
		{
			name: "no matched transactions",
			matchedTransactions: func(t *testing.T) MatchedTransactions {
				return []journal.Transaction{}
			},
			inputPostings: func() []journal.Posting { return []journal.Posting{} },
			success:       false,
		},
		{
			name: "one single matched transaction",
			matchedTransactions: func(*testing.T) MatchedTransactions {
				return []journal.Transaction{*tu.Transaction_1(t)}
			},
			inputPostings: func() []journal.Posting { return []journal.Posting{} },
			success:       true,
			expected:      "ACC1",
		},
		{
			name: "two matched transactions",
			matchedTransactions: func(*testing.T) MatchedTransactions {
				return []journal.Transaction{*tu.Transaction_1(t), *tu.Transaction_2(t)}
			},
			inputPostings: func() []journal.Posting { return []journal.Posting{} },
			success:       true,
			expected:      "ACC1",
		},
		{
			name: "with previous input posting",
			matchedTransactions: func(*testing.T) MatchedTransactions {
				return []journal.Transaction{*tu.Transaction_1(t)}
			},
			inputPostings: func() []journal.Posting {
				return []journal.Posting{
					{
						Account: "ACC1",
						Ammount: finance.Ammount{},
					},
				}
			},
			success:  true,
			expected: "ACC2",
		},
	}

	// Default setup if test does not define one
	defaultSetup := func(t *testing.T, c *testcontext) {}

	// Run test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			if tc.setup != nil {
				tc.setup(t, c)
			} else {
				defaultSetup(t, c)
			}
			c.accountguesser, err = NewMatchedTransactionsAccountGuesser()
			if err != nil {
				t.Fatal(err)
			}
			c.accountguesser.SetMatchedTransactions(tc.matchedTransactions(t))
			c.accountguesser.SetInputPostings(tc.inputPostings())
			actual, success := c.accountguesser.Guess()
			assert.Equal(t, tc.success, success)
			if tc.success {
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestLastTransactionAccountGuesser(t *testing.T) {
	type testcontext struct {
		accountguesser *LastTransactionAccountGuesser
	}
	type testcase struct {
		name               string
		transactionHistory func() MatchedTransactions
		success            bool
		expected           journal.Account
	}
	var testcases = []testcase{
		{
			name: "no matched transactions",
			transactionHistory: func() MatchedTransactions {
				return []journal.Transaction{}
			},
			success: false,
		},
		{
			name: "uses last transaction first posting",
			transactionHistory: func() MatchedTransactions {
				return []journal.Transaction{
					{
						Posting: []journal.Posting{
							{
								Account: "supermarket",
							},
						},
					},
					{
						Posting: []journal.Posting{
							{
								Account: "savings",
							},
						},
					},
				}
			},
			success:  true,
			expected: "savings",
		},
	}

	// Run test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			c := new(testcontext)
			c.accountguesser, err = NewLastTransactionAccountGuesser()
			if err != nil {
				t.Fatal(err)
			}
			c.accountguesser.SetTransactionHistory(tc.transactionHistory())
			actual, success := c.accountguesser.Guess()
			assert.Equal(t, tc.success, success)
			if tc.success {
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestCompositeAccountGuesser(t *testing.T) {
	type testcontext struct {
		accountguesser *CompositeAccountGuesser
	}
	type testcase struct {
		name             string
		composedGuessers func(ctrl *gomock.Controller) []IAccountGuesser
		success          bool
		expected         journal.Account
	}
	var testcases = []testcase{
		{
			name: "no composed guessers",
			composedGuessers: func(ctrl *gomock.Controller) []IAccountGuesser {
				return []IAccountGuesser{}
			},
			success: false,
		},
		{
			name: "single composed guesser (succcess)",
			composedGuessers: func(ctrl *gomock.Controller) []IAccountGuesser {
				accountGuesser := NewMockIAccountGuesser(ctrl)
				accountGuesser.EXPECT().Guess().Return(journal.Account("savings"), true)
				return []IAccountGuesser{accountGuesser}
			},
			success:  true,
			expected: "savings",
		},
		{
			name: "single composed guesser (failure)",
			composedGuessers: func(ctrl *gomock.Controller) []IAccountGuesser {
				accountGuesser := NewMockIAccountGuesser(ctrl)
				accountGuesser.EXPECT().Guess().Return(journal.Account(""), false)
				return []IAccountGuesser{accountGuesser}
			},
			success: false,
		},
		{
			name: "two composed guesser (first success)",
			composedGuessers: func(ctrl *gomock.Controller) []IAccountGuesser {
				accountGuesserOne := NewMockIAccountGuesser(ctrl)
				accountGuesserOne.EXPECT().Guess().Return(journal.Account("savings1"), true)
				accountGuesserTwo := NewMockIAccountGuesser(ctrl)
				return []IAccountGuesser{accountGuesserOne, accountGuesserTwo}
			},
			success:  true,
			expected: "savings1",
		},
		{
			name: "two composed guesser (second success)",
			composedGuessers: func(ctrl *gomock.Controller) []IAccountGuesser {
				accountGuesserOne := NewMockIAccountGuesser(ctrl)
				accountGuesserOne.EXPECT().Guess().Return(journal.Account(""), false)
				accountGuesserTwo := NewMockIAccountGuesser(ctrl)
				accountGuesserTwo.EXPECT().Guess().Return(journal.Account("savings2"), true)
				return []IAccountGuesser{accountGuesserOne, accountGuesserTwo}
			},
			success:  true,
			expected: "savings2",
		},
		{
			name: "two composed guesser (failure)",
			composedGuessers: func(ctrl *gomock.Controller) []IAccountGuesser {
				accountGuesserOne := NewMockIAccountGuesser(ctrl)
				accountGuesserOne.EXPECT().Guess().Return(journal.Account(""), false)
				accountGuesserTwo := NewMockIAccountGuesser(ctrl)
				accountGuesserTwo.EXPECT().Guess().Return(journal.Account(""), false)
				return []IAccountGuesser{accountGuesserOne, accountGuesserTwo}
			},
			success: false,
		},
	}

	// Run test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			var err error
			c := new(testcontext)
			composedGuessers := tc.composedGuessers(ctrl)
			c.accountguesser, err = NewCompositeAccountGuesser(composedGuessers...)
			if err != nil {
				t.Fatal(err)
			}
			actual, success := c.accountguesser.Guess()
			assert.Equal(t, tc.success, success)
			if tc.success {
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}
