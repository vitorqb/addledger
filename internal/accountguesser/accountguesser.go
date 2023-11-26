// accountguesser is a package that provides a way to guess which account the user
// wants to input.
package accountguesser

import (
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/statementloader"
)

//go:generate $MOCKGEN --source=accountguesser.go --destination=../../mocks/accountguesser/accountguesser_mock.go

// MatchedTransactions represents a collection of matched transactions.
type MatchedTransactions []journal.Transaction

// TransactionHistory represents a collection of transactions already present in the journal.
type TransactionHistory []journal.Transaction

// Inputs are the inputs that are used for guessing.
type Inputs struct {
	// MatchingTransactions are the transactions that match the current user input.
	MatchingTransactions MatchedTransactions
	// PostingInputs are the postings that have been inputted so far by the user,
	// in the current journal entry.
	PostingInputs []journal.Posting
	// Description is the description the user has inputted in the description field.
	Description string
	// TransactionHistory is the transaction history in the journal.
	TransactionHistory TransactionHistory
	// StatementEntry is the statement entry that has been loaded and is being used
	// for the current journal entry.
	StatementEntry statementloader.StatementEntry
}

// AccountGuesser is a strategy for guessing the account an user may want for an journal entry.
type AccountGuesser interface {
	Guess(inputs Inputs) (guess journal.Account, success bool)
}

var _ AccountGuesser = &MatchedTransactionsGuesser{}

// MatchedTransactionsGuesser uses the matched transactions and
// returns the suggested account based on the best-matched
// transaction. `inputPostings` must be the postings the user has
// already entered for the new transaction.
type MatchedTransactionsGuesser struct{}

// Guess implements IAccountGuesser.
func (*MatchedTransactionsGuesser) Guess(inputs Inputs) (acc journal.Account, success bool) {
	if len(inputs.MatchingTransactions) == 0 {
		return "", false
	}
	matchedTransaction := inputs.MatchingTransactions[0]

	// We had a match, so find posting the user is entering
	desiredPostingIndex := len(inputs.PostingInputs)

	// If the user has already entered more posting than matched transaction,
	// we can't use it.
	if desiredPostingIndex >= len(matchedTransaction.Posting) {
		return "", false
	}

	// Otherwise get the account from the posting with same index.
	matchedPosting := matchedTransaction.Posting[desiredPostingIndex]
	return journal.Account(matchedPosting.Account), true
}

// NewMatchedTransactionsAccountGuesser returns a new implementation of MatchedTransactionsGuesser
func NewMatchedTransactionsAccountGuesser() (*MatchedTransactionsGuesser, error) {
	return &MatchedTransactionsGuesser{}, nil
}

// LastTransactionAccountGuesser uses the last entered transaction to try
// to guess an account.
type LastTransactionAccountGuesser struct{}

var _ AccountGuesser = &LastTransactionAccountGuesser{}

func (ag *LastTransactionAccountGuesser) Guess(inputs Inputs) (acc journal.Account, success bool) {
	historyLen := len(inputs.TransactionHistory)
	if historyLen == 0 {
		return "", false
	}
	lastTransaction := inputs.TransactionHistory[historyLen-1]
	if len(lastTransaction.Posting) == 0 {
		return "", false
	}
	firstPosting := lastTransaction.Posting[0]
	return journal.Account(firstPosting.Account), true
}

func NewLastTransactionAccountGuesser() (*LastTransactionAccountGuesser, error) {
	return &LastTransactionAccountGuesser{}, nil
}

// StatementAccountGuesser uses the current statement entry to try
// to guess an account
type StatementAccountGuesser struct{}

var _ AccountGuesser = &StatementAccountGuesser{}

// Guess tries to guess an account based on a loaded statement entry. If there
// is a statement entry with an acconut that does not yet exist in the
// input postings, it returns that account.
func (ag *StatementAccountGuesser) Guess(inputs Inputs) (acc journal.Account, success bool) {
	if inputs.StatementEntry.Account == "" {
		return "", false
	}
	for _, posting := range inputs.PostingInputs {
		if posting.Account == inputs.StatementEntry.Account {
			return "", false
		}
	}
	return journal.Account(inputs.StatementEntry.Account), true
}

func NewStatementAccountGuesser() (*StatementAccountGuesser, error) {
	return &StatementAccountGuesser{}, nil
}

// CompositeAccountGuesser composes N different account guessers
type CompositeAccountGuesser struct {
	composedGuessers []AccountGuesser
}

var _ AccountGuesser = &CompositeAccountGuesser{}

// Guess implements IAccountGuesser.
func (ag *CompositeAccountGuesser) Guess(inputs Inputs) (guess journal.Account, success bool) {
	for _, composedGuesser := range ag.composedGuessers {
		if guess, ok := composedGuesser.Guess(inputs); ok {
			return guess, ok
		}
	}
	return journal.Account(""), false
}

func NewCompositeAccountGuesser(accGuessers ...AccountGuesser) (*CompositeAccountGuesser, error) {
	return &CompositeAccountGuesser{accGuessers}, nil
}
